package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"github.com/zigzter/chatterm/utils"
)

var (
	msgRegex *regexp.Regexp = regexp.MustCompile(`\bPRIVMSG\b`)
	cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)
)

func establishWSConnection(channel string, username string, oath string) {
	socketUrl := "ws://irc-ws.chat.twitch.tv:80"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to socket server:", err)
	}
	defer conn.Close()
	err1 := conn.WriteMessage(websocket.TextMessage, []byte("CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands"))
	if err1 != nil {
		log.Println("Error writing message:", err1)
	}
	setPass := fmt.Sprintf("PASS %s", oath)
	log.Println(setPass)
	conn.WriteMessage(websocket.TextMessage, []byte(setPass))
	setUser := fmt.Sprintf("NICK %s", username)
	conn.WriteMessage(websocket.TextMessage, []byte(setUser))
	joinChannel := fmt.Sprintf("JOIN #%s", channel)
	conn.WriteMessage(websocket.TextMessage, []byte(joinChannel))
	for {
		// Read messages and handle them
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			return
		}
		if messageType == websocket.TextMessage {
			rawIrcMessage := strings.TrimSpace(string(message))
			if msgRegex.MatchString(rawIrcMessage) {
				chatMessage := utils.MessageParser(rawIrcMessage)
				fmt.Printf("%s: %s \n", chatMessage.DisplayName, chatMessage.Message)
			} else {
				fmt.Println(rawIrcMessage)
			}
		}
	}
}

var (
	Channel  string
	Username string
	Oauth    string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects to a Twitch chat",
	Long:  `Connects to a Twitch chat`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running cmd", Channel, Username, Oauth)
		go establishWSConnection(Channel, Username, Oauth)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	},
}

func init() {
	connectCmd.Flags().StringVarP(&Channel, "channel", "c", "", "The Twitch channel to join")
	connectCmd.Flags().StringVarP(&Username, "username", "u", "", "Your username on Twitch")
	connectCmd.Flags().StringVarP(&Oauth, "oauth", "o", "", "The Oath string, in format oauth:xyz123")
	connectCmd.MarkFlagRequired("channel")
	connectCmd.MarkFlagRequired("username")
	connectCmd.MarkFlagRequired("oauth")
	connectCmd.Println(Channel, Username, Oauth)

	rootCmd.AddCommand(connectCmd)
}
