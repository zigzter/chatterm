package twitch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/types"
)

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

func SendTwitchCommand(command types.TwitchCommand, args string) error {
	cmdDetails := RequestMap[command]
	argParts := strings.Split(args, " ")
	targetUser := string(argParts[0])
	duration := string(argParts[1])
	channel := viper.GetString("channel")
	moderatorId := viper.GetString("userId")
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
	url := rootUrl + cmdDetails.Endpoint + "?broadcaster_id=" + channel + "&moderator_id=" + moderatorId
	token := viper.GetString("token")
	requestBody, err := json.Marshal(map[string]map[string]string{
		"data": {"user_id": userId, "duration": duration},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(cmdDetails.Method, url, bytes.NewBuffer(requestBody))
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
