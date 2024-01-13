package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
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

// SendTwitchCommand sends a request to the Twitch Helix API to perform a command
func SendTwitchCommand(command types.TwitchCommand, args []string) (interface{}, error) {
	req, err := ConstructRequest(command, args)
	if err != nil {
		return nil, err
	}
	client := httpClient()
	resp, err := client.Do(req)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result interface{}
	switch command {
	case types.Info:
		var userData types.UserData
		if err := json.Unmarshal(bodyBytes, &userData); err != nil {
			return nil, err
		}
		result = userData
	}
	return result, nil
}

// ConstructRequest creates the specific request to be sent by SendTwitchCommand
func ConstructRequest(command types.TwitchCommand, args []string) (*http.Request, error) {
	cmdDetails := RequestMap[command]
	rootUrl := "https://api.twitch.tv/helix"
	url := rootUrl + cmdDetails.Endpoint
	var req *http.Request
	var err error
	switch command {
	case types.Info:
		req, err = http.NewRequest(cmdDetails.Method, url, nil)
		query := req.URL.Query()
		query.Add("login", args[0])
		req.URL.RawQuery = query.Encode()
	case types.Ban:
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
			data, err := SendTwitchCommand(types.Info, args)
			if err != nil {
				return nil, err
			}
			if user, ok := data.(types.UserData); ok {
				userId = user.ID
				db.InsertUserMap(sql, targetUser, userId)
			}
		}
		requestBody, err := json.Marshal(map[string]map[string]string{
			"data": {"user_id": userId, "duration": duration},
		})
		req, err = http.NewRequest(cmdDetails.Method, url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
	case types.Unban:
		targetUser := string(args[0])
		sql := db.OpenDB()
		userId, err := db.GetUserId(sql, targetUser)
		if err != nil {
			return nil, err
		}
		if userId == "" {
			data, err := SendTwitchCommand(types.Info, args)
			if err != nil {
				return nil, err
			}
			if user, ok := data.(types.UserData); ok {
				userId = user.ID
				db.InsertUserMap(sql, targetUser, userId)
			}
		}
		req, err = http.NewRequest(cmdDetails.Method, url, nil)
		q := req.URL.Query()
		q.Add("user_id", userId)
		req.URL.RawQuery = q.Encode()
	case types.Clear:
		req, err = http.NewRequest(cmdDetails.Method, url, nil)
	case types.Delete:
		if len(args) < 1 {
			return nil, errors.New("Please provide the id of the message to delete")
		}
		req, err = http.NewRequest(cmdDetails.Method, url, nil)
		req.URL.Query().Add("message_id", args[0])
	default:
		return nil, errors.New("Invalid command")
	}
	channelid := viper.GetString("channelid")
	moderatorId := viper.GetString("userid")
	token := viper.GetString("token")
	query := req.URL.Query()
	query.Add("broadcaster_id", channelid)
	query.Add("moderator_id", moderatorId)
	req.URL.RawQuery = query.Encode()
	log.Println(req.URL.RawQuery)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", ClientId)
	return req, err
}
