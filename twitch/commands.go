package twitch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/types"
)

// FetchUser retrieves the account data of the provided username from the Twitch API
func FetchUser(username string) (types.UserData, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	query := req.URL.Query()
	query.Add("login", username)
	req.URL.RawQuery = query.Encode()
	token := viper.GetString("token")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", ClientId)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching user:", err)
	}
	defer resp.Body.Close()
	var response types.FetchUserResp
	jsonErr := json.NewDecoder(resp.Body).Decode(&response)
	if jsonErr != nil {
		log.Println("Json error:", jsonErr)
	}
	user := response.Data[0]
	return user, nil
}

// SendTwitchCommand sends a request to the Twitch Helix API to perform a command
func SendTwitchCommand(command types.TwitchCommand, args []string) error {
	cmdDetails := RequestMap[command]
	targetUser := string(args[0])
	duration := "0"
	if len(args) > 1 {
		duration = string(args[1])
	}
	channelid := viper.GetString("channelid")
	moderatorId := viper.GetString("userid")
	sql := db.OpenDB()
	userId, err := db.GetUserId(sql, targetUser)
	if userId == "" {
		user, err := FetchUser(targetUser)
		if err != nil {
			return err
		}
		userId = user.ID
		db.InsertUserMap(sql, targetUser, userId)
	}
	rootUrl := "https://api.twitch.tv/helix"
	url := rootUrl + cmdDetails.Endpoint
	token := viper.GetString("token")
	requestBody, err := json.Marshal(map[string]map[string]string{
		"data": {"user_id": userId, "duration": duration},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(cmdDetails.Method, url, bytes.NewBuffer(requestBody))
	query := req.URL.Query()
	query.Add("broadcaster_id", channelid)
	query.Add("moderator_id", moderatorId)
	req.URL.RawQuery = query.Encode()
	log.Println("sending to url:", req.URL)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", ClientId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
