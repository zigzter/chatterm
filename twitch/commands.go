package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
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

func fireRequest(req *http.Request) ([]byte, error) {
	client := httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	if resp == nil {
		return nil, errors.New("HTTP request returned nil response")
	}
	defer resp.Body.Close()

	var buff bytes.Buffer
	_, err = io.Copy(&buff, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, errors.New(string(buff.Bytes()))
	}
	return buff.Bytes(), nil
}

type APIRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Query   map[string]string
	Body    []byte
}

func NewAPIRequest(method, url string, body []byte) *APIRequest {
	channelID := viper.GetString(utils.ChannelIDKey)
	moderatorID := viper.GetString(utils.UserIDKey)
	token := viper.GetString(utils.TokenKey)
	return &APIRequest{
		Method: method,
		URL:    url,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
			"Client-Id":     ClientId,
		},
		Query: map[string]string{
			"moderator_id":   moderatorID,
			"broadcaster_id": channelID,
		},
		Body: body,
	}
}

func (req *APIRequest) AddHeader(key, value string) *APIRequest {
	req.Headers[key] = value
	return req
}

func (req *APIRequest) AddQuery(key, value string) *APIRequest {
	req.Query[key] = value
	return req
}

func (req *APIRequest) Execute() ([]byte, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	query := httpReq.URL.Query()
	for key, value := range req.Query {
		query.Add(key, value)
	}
	httpReq.URL.RawQuery = query.Encode()
	return fireRequest(httpReq)
}

func isValidCommand(command string) bool {
	switch types.TwitchCommand(command) {
	case types.Ban, types.Clear, types.Unban, types.Delete,
		types.Info, types.Shield, types.Announce,
		types.FollowersOnly, types.SubOnly, types.Slow,
		types.EmoteOnly, types.Shoutout, types.Warn:
		return true
	}
	return false
}

// SendTwitchCommand sends a request to the Twitch Helix API to perform a command
func SendTwitchCommand(command types.TwitchCommand, args []string) (interface{}, error) {
	if isValid := isValidCommand(string(command)); !isValid {
		return nil, errors.New("Invalid command")
	}
	switch command {
	case types.Ban:
		return sendBanRequest(args)
	case types.Unban:
		return sendUnbanRequest(args)
	case types.Info:
		return sendInfoRequest(args[0])
	case types.Clear:
		return sendClearRequest()
	case types.Shoutout:
		return sendShoutout(args[0])
	case types.Delete:
		return sendDeleteRequest(args)
	case types.Shield:
		return sendShieldRequest(args)
	case types.Warn:
		return sendWarning(args[0], strings.Join(args[1:], " "))
	case types.Announce:
		return sendAnnouncement(strings.Join(args, " "))
	case types.Slow, types.SubOnly, types.FollowersOnly:
		var duration string
		if len(args) > 0 {
			duration = args[0]
		}
		return sendUpdateChatRequest(command, duration)
	}
	return nil, fmt.Errorf("Unknown command: %s", command)
}

func sendUpdateChatRequest(mode types.TwitchCommand, duration string) (*types.UpdateChatSettingsData, error) {
	bodyData := make(map[string]interface{})
	shouldEnable := duration != "off"
	intDuration, _ := strconv.Atoi(duration)
	switch mode {
	case types.EmoteOnly:
		bodyData["emote_only"] = shouldEnable
	case types.FollowersOnly:
		bodyData["follower_mode"] = shouldEnable
		if intDuration > 0 {
			bodyData["follower_mode_duration"] = intDuration
		}
	case types.SubOnly:
		bodyData["subscriber_mode"] = shouldEnable
	case types.Slow:
		bodyData["slow_mod"] = shouldEnable
		if intDuration > 0 {
			bodyData["follower_mode_duration"] = intDuration
		}
	default:
		return nil, errors.New("Invalid chat setting")
	}
	cmdDetails := RequestMap[types.TwitchCommand(mode)]
	url := rootUrl + cmdDetails.Endpoint
	requestBody, err := json.Marshal(map[string]interface{}{
		"data": bodyData,
	})
	response, err := NewAPIRequest(cmdDetails.Method, url, requestBody).
		AddHeader("Content-Type", "application/json").
		Execute()
	if err != nil {
		return nil, err
	}
	var updateChatSettingsData types.UpdateChatSettingsData
	json.Unmarshal(response, &updateChatSettingsData)
	return &updateChatSettingsData, nil
}

func SendUserRequest(username string) (*types.UserResp, error) {
	cmdDetails := RequestMap[types.User]
	url := rootUrl + cmdDetails.Endpoint
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("login", username).
		Execute()
	if err != nil {
		return nil, err
	}
	var userData types.UserResp
	json.Unmarshal(response, &userData)
	return &userData, nil
}

func sendFollowersRequest(args []string) (*types.FollowersResp, error) {
	cmdDetails := RequestMap[types.GetFollowers]
	url := rootUrl + cmdDetails.Endpoint
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("user_id", args[0]).
		Execute()
	if err != nil {
		return nil, err
	}
	var followerData types.FollowersResp
	json.Unmarshal(response, &followerData)
	return &followerData, nil
}

// getUserColor retrieves the target user's chat color from the Twitch API
func getUserColor(userId string) (*types.ColorResp, error) {
	cmdDetails := RequestMap[types.Color]
	url := rootUrl + cmdDetails.Endpoint
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("user_id", userId).
		Execute()
	if err != nil {
		return nil, err
	}
	var colorData types.ColorResp
	json.Unmarshal(response, &colorData)
	return &colorData, nil
}

// sendInfoRequest hits several API endpoints and returns a single collection of user data
func sendInfoRequest(username string) (*types.UserInfo, error) {
	userInfo := types.UserInfo{}
	userResp, err := SendUserRequest(username)
	if err != nil {
		return nil, err
	}
	if len(userResp.Data) == 0 {
		return nil, errors.New("User does not exist")
	}
	userInfo.Details = userResp.Data[0]
	colorData, err := getUserColor(userInfo.Details.ID)
	if err != nil {
		log.Println("Error retrieving color of:", userInfo.Details.DisplayName)
	}
	userInfo.Color = colorData.Data[0].Color
	userType := viper.GetString(utils.ChannelUserTypeKey)
	hasModPrivs := userType == "moderator" || userType == "broadcaster"
	if !hasModPrivs {
		// The follower endpoint requires moderator privileges
		return &userInfo, nil
	}
	followerResp, err := sendFollowersRequest([]string{userInfo.Details.ID})
	if err != nil {
		return nil, err
	}
	if len(followerResp.Data) > 0 {
		userInfo.Following = followerResp.Data[0]
	}
	return &userInfo, nil
}

func SendLiveChannelsRequest(userID string) (*types.LiveChannelsResp, error) {
	cmdDetails := RequestMap[types.LiveChannels]
	url := rootUrl + cmdDetails.Endpoint
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("user_id", userID).
		Execute()
	if err != nil {
		return nil, err
	}
	var liveChannels types.LiveChannelsResp
	json.Unmarshal(response, &liveChannels)
	return &liveChannels, nil
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
		user, err := SendUserRequest(args[0])
		if err != nil {
			return nil, err
		}
		userId = user.Data[0].ID
		db.InsertUserMap(sql, targetUser, userId)
	}
	requestBody, err := json.Marshal(map[string]map[string]string{
		"data": {"user_id": userId, "duration": duration},
	})
	response, err := NewAPIRequest(cmdDetails.Method, url, requestBody).
		AddHeader("Content-Type", "application/json").
		Execute()
	if err != nil {
		return nil, err
	}
	var banResponse types.UserBanResp
	json.Unmarshal(response, &banResponse)
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
		user, err := SendUserRequest(args[0])
		if err != nil {
			return nil, err
		}
		userId = user.Data[0].ID
		db.InsertUserMap(sql, targetUser, userId)
	}
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("user_id", userId).
		Execute()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func sendClearRequest() (any, error) {
	cmdDetails := RequestMap[types.Clear]
	url := rootUrl + cmdDetails.Endpoint
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		Execute()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func sendDeleteRequest(args []string) (any, error) {
	cmdDetails := RequestMap[types.Delete]
	url := rootUrl + cmdDetails.Endpoint
	if len(args) < 1 {
		return nil, errors.New("Please provide the id of the message to delete")
	}
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("message_id", args[0]).
		Execute()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func sendShieldRequest(args []string) (*types.ShieldResp, error) {
	if len(args) == 0 {
		return nil, errors.New("Please provide 'on' or 'off' as an argument")
	}
	status := args[0]
	if status != "on" && status != "off" {
		return nil, errors.New("Please provide 'on' or 'off' as an argument")
	}
	cmdDetails := RequestMap[types.Shield]
	url := rootUrl + cmdDetails.Endpoint
	enable := true
	if status == "off" {
		enable = false
	}
	requestBody, err := json.Marshal(map[string]bool{"is_active": enable})
	response, err := NewAPIRequest(cmdDetails.Method, url, requestBody).
		AddHeader("Content-Type", "application/json").
		Execute()
	if err != nil {
		return nil, err
	}
	var shieldResponse types.ShieldResp
	json.Unmarshal(response, &shieldResponse)
	return &shieldResponse, nil
}

func sendShoutout(username string) (any, error) {
	sql := db.OpenDB()
	targetUserId, err := db.GetUserId(sql, username)
	if err != nil {
		return nil, err
	}
	if targetUserId == "" {
		user, err := SendUserRequest(username)
		if err != nil {
			return nil, err
		}
		targetUserId = user.Data[0].ID
		db.InsertUserMap(sql, username, targetUserId)
	}
	originUserId := viper.GetString(utils.ChannelIDKey)
	cmdDetails := RequestMap[types.Shoutout]
	url := rootUrl + cmdDetails.Endpoint
	response, err := NewAPIRequest(cmdDetails.Method, url, nil).
		AddQuery("from_broadcaster_id", originUserId).
		AddQuery("to_broadcaster_id", targetUserId).
		Execute()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func sendWarning(username, reason string) (*types.WarnResp, error) {
	sql := db.OpenDB()
	targetUserId, err := db.GetUserId(sql, username)
	if err != nil {
		return nil, err
	}
	if targetUserId == "" {
		user, err := SendUserRequest(username)
		if err != nil {
			return nil, err
		}
		targetUserId = user.Data[0].ID
		db.InsertUserMap(sql, username, targetUserId)
	}
	cmdDetails := RequestMap[types.Warn]
	url := rootUrl + cmdDetails.Endpoint
	requestBody, err := json.Marshal(map[string]map[string]string{
		"data": {"user_id": targetUserId, "reason": reason},
	})
	response, err := NewAPIRequest(cmdDetails.Method, url, requestBody).
		AddHeader("Content-Type", "application/json").
		Execute()
	if err != nil {
		return nil, err
	}
	var warnResponse types.WarnResp
	json.Unmarshal(response, &warnResponse)
	return &warnResponse, nil
}

func sendAnnouncement(text string) (any, error) {
	cmdDetails := RequestMap[types.Announce]
	url := rootUrl + cmdDetails.Endpoint
	requestBody, err := json.Marshal(map[string]string{
		"message": text,
	})
	if err != nil {
		return nil, err
	}
	_, err = NewAPIRequest(cmdDetails.Method, url, requestBody).
		AddHeader("Content-Type", "application/json").
		Execute()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
