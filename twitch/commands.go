package twitch

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

type TwitchCommand string

const (
	Ban   TwitchCommand = "ban"
	Clear TwitchCommand = "clear"
)

func SendTwitchCommand(command TwitchCommand) error {
	url := "https://api.twitch.tv/helix/"
	token := viper.GetString("token")
	clientId := viper.GetString("clientid")
	requestBody, err := json.Marshal(map[string]TwitchCommand{"data": command})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", clientId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
