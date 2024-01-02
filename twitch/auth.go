package twitch

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

var serverAddr = "localhost:3030"

func startLocalServer(ready chan<- struct{}) {
	http.HandleFunc("/token/", func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.Path, "/token/")
		if token != "" {
			utils.SaveConfig(map[string]interface{}{
				"token": token,
			})
			fmt.Fprintln(w, "Token received, you can close this window.")
		} else {
			fmt.Fprintln(w, "Failed to retrieve token.")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := "./twitch/token.html"
		http.ServeFile(w, r, filePath)
	})

	httpServer := &http.Server{Addr: serverAddr}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	ready <- struct{}{}
}

func PromptTwitchAuth() {
	scopes := []string{
		"chat:read",
		"chat:edit",
		"user:read:chat",
		"channel:moderate",
	}
	clientId := "x6pl99d1tq9mqys6y2bmr59ahw9nik"
	redirectUrl := serverAddr
	scope := strings.Join(scopes, " ")
	state := utils.GenerateRandomString(10)
	url := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?response_type=token&client_id=%s&redirect_uri=http://%s&scope=%s&state=%s",
		clientId,
		redirectUrl,
		scope,
		state,
	)
	utils.OpenBrowser(url)
}

func StartAuthenticationProcess() tea.Cmd {
	return func() tea.Msg {
		ready := make(chan struct{})
		go startLocalServer(ready)
		<-ready
		PromptTwitchAuth()
		return types.AuthResultMsg{}
	}
}
