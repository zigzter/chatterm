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

func StartLocalServer(ready chan<- struct{}, externalMsgs chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		http.HandleFunc("/token/", func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.URL.Path, "/token/")
			if token != "" {
				utils.SaveConfig(map[string]interface{}{
					"token": token,
				})
				fmt.Fprintln(w, "Token received, you can close this window.")
				externalMsgs <- types.TokenReceivedMsg{}
			} else {
				fmt.Fprintln(w, "Failed to retrieve token.")
				externalMsgs <- types.TokenReceivedMsg{}
			}
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			filePath := "./twitch/token.html"
			http.ServeFile(w, r, filePath)
		})

		httpServer := &http.Server{Addr: serverAddr}
		go func() {
			ready <- struct{}{}
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe error: %v", err)
			}
		}()
		return types.ServerStartedMsg{}
	}
}

func PromptTwitchAuth() tea.Cmd {
	return func() tea.Msg {
		scopes := []string{
			"chat:read",
			"chat:edit",
			"user:read:chat",
			"channel:moderate",
			"moderator:manage:chat_messages",
			"moderator:manage:banned_users",
		}
		redirectUrl := serverAddr
		scope := strings.Join(scopes, " ")
		state := utils.GenerateRandomString(10)
		url := fmt.Sprintf(
			"https://id.twitch.tv/oauth2/authorize?response_type=token&client_id=%s&redirect_uri=http://%s&scope=%s&state=%s",
			ClientId,
			redirectUrl,
			scope,
			state,
		)
		utils.OpenBrowser(url)
		return types.AuthOpenedMsg{}
	}
}
