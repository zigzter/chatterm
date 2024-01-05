package twitch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/types"
)

func SendTwitchCommand(command types.TwitchCommand, args string) error {
	cmdDetails := RequestMap[command]
	argParts := strings.SplitN(args, " ", 2)
	channel := viper.GetString("channel")
	username := viper.GetString("username")
	rootUrl := "https://api.twitch.tv/helix"
	url := rootUrl + cmdDetails.Endpoint + "?broadcaster_id=" + channel + "&moderator_id=" + username
	token := viper.GetString("token")
	requestBody, err := json.Marshal(map[string]map[string]string{
		"data": {"user_id": argParts[0], "duration": argParts[1]},
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
