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

func establishWSConnection(channel string) {
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
	conn.WriteMessage(websocket.TextMessage, []byte("PASS replaceme"))
	conn.WriteMessage(websocket.TextMessage, []byte("NICK zigzter"))
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("#%s", channel)))
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
				fmt.Printf("%+v\n", chatMessage)
			}
		}
	}
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects to a Twitch chat",
	Long:  `Connects to a Twitch chat`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect called!!!", args)

		go establishWSConnection(args[0])

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		// if joinErr != nil {
		// 	log.Fatal("Error joining room:", joinErr)
		// }
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
