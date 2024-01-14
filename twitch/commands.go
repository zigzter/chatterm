package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/types"
)

var (
	clientInstance *http.Client
	once           sync.Once
)

const rootUrl string = "https://api.twitch.tv/helix"

// httpClient creates a new http client, reusing an existing instance if it exists.
func httpClient() *http.Client {
	once.Do(func() {
		clientInstance = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
			Timeout: 10 * time.Second,
		}
	})
	return clientInstance
}

func augmentRequest(req *http.Request) *http.Request {
	channelid := viper.GetString("channelid")
	moderatorId := viper.GetString("userid")
	token := viper.GetString("token")
	query := req.URL.Query()
	query.Add("broadcaster_id", channelid)
	query.Add("moderator_id", moderatorId)
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", ClientId)
	return req
}

func fireRequest(req *http.Request) (string, []byte, error) {
	client := httpClient()
	resp, err := client.Do(req)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	return resp.Status, bodyBytes, nil
}

// SendTwitchCommand sends a request to the Twitch Helix API to perform a command
func SendTwitchCommand(command types.TwitchCommand, args []string) (interface{}, error) {
	switch command {
	case types.Ban:
		return sendBanRequest(args)
	case types.Unban:
		return sendUnbanRequest(args)
	case types.Info:
		return sendInfoRequest(args)
	case types.Clear:
		return sendClearRequest()
	case types.Delete:
		return sendDeleteRequest(args)
	}
	return nil, fmt.Errorf("Unknown command: %s", command)
}

func sendInfoRequest(args []string) (*types.UserData, error) {
	cmdDetails := RequestMap[types.Info]
	url := rootUrl + cmdDetails.Endpoint
	req, err := http.NewRequest(cmdDetails.Method, url, nil)
	if err != nil {
		return nil, err
	}
	req = augmentRequest(req)
	query := req.URL.Query()
	query.Add("login", args[0])
	req.URL.RawQuery = query.Encode()
	_, bytes, err := fireRequest(req)
	if err != nil {
		return nil, err
	}
	var userData types.UserData
	json.Unmarshal(bytes, &userData)
	return &userData, nil
}

func sendBanRequest(args []string) (*types.UserBanResp, error) {
	cmdDetails := RequestMap[types.Ban]
	url := rootUrl + cmdDetails.Endpoint
	targetUser := string(args[0])
	duration := "0"
	if len(args) > 1 {
		duration = string(args[1])
	}
	sql := db.OpenDB()
	userId, err := db.GetUserId(sql, targetUser)
	if err != nil {
		return nil, err
	}
	if userId == "" {
		user, err := sendInfoRequest(args)
		if err != nil {
			return nil, err
		}
		userId = user.Data[0].ID
		db.InsertUserMap(sql, targetUser, userId)
	}
	requestBody, err := json.Marshal(map[string]map[string]string{
		"data": {"user_id": userId, "duration": duration},
	})
	req, err := http.NewRequest(cmdDetails.Method, url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req = augmentRequest(req)
	_, bytes, err := fireRequest(req)
	if err != nil {
		return nil, err
	}
	var banResponse types.UserBanResp
	json.Unmarshal(bytes, &banResponse)
	return &banResponse, nil
}

// TODO: this any business to match return signatures seems silly, find a better way
func sendUnbanRequest(args []string) (any, error) {
	cmdDetails := RequestMap[types.Unban]
	url := rootUrl + cmdDetails.Endpoint
	targetUser := string(args[0])
	sql := db.OpenDB()
	userId, err := db.GetUserId(sql, targetUser)
	if err != nil {
		return nil, err
	}
	if userId == "" {
		user, err := sendInfoRequest(args)
		if err != nil {
			return nil, err
		}
		userId = user.Data[0].ID
		db.InsertUserMap(sql, targetUser, userId)
	}
	req, err := http.NewRequest(cmdDetails.Method, url, nil)
	q := req.URL.Query()
	q.Add("user_id", userId)
	req.URL.RawQuery = q.Encode()
	req = augmentRequest(req)
	status, bytes, err := fireRequest(req)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(status, "20") {
		var apiErr types.TwitchAPIError
		json.Unmarshal(bytes, &apiErr)
		return nil, apiErr
	}
	return nil, nil
}

func sendClearRequest() (any, error) {
	cmdDetails := RequestMap[types.Clear]
	url := rootUrl + cmdDetails.Endpoint
	req, err := http.NewRequest(cmdDetails.Method, url, nil)
	if err != nil {
		return nil, err
	}
	req = augmentRequest(req)
	status, bytes, err := fireRequest(req)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(status, "20") {
		var apiErr types.TwitchAPIError
		json.Unmarshal(bytes, &apiErr)
		return nil, apiErr
	}
	return nil, nil
}

func sendDeleteRequest(args []string) (any, error) {
	cmdDetails := RequestMap[types.Delete]
	url := rootUrl + cmdDetails.Endpoint
	if len(args) < 1 {
		return nil, errors.New("Please provide the id of the message to delete")
	}
	req, err := http.NewRequest(cmdDetails.Method, url, nil)
	q := req.URL.Query()
	q.Add("message_id", args[0])
	req.URL.RawQuery = q.Encode()
	req = augmentRequest(req)
	status, bytes, err := fireRequest(req)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(status, "20") {
		var apiErr types.TwitchAPIError
		json.Unmarshal(bytes, &apiErr)
		return nil, apiErr
	}
	return nil, nil
}
