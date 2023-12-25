package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	msgRegex  *regexp.Regexp = regexp.MustCompile(`\bPRIVMSG\b`)
	subRegex  *regexp.Regexp = regexp.MustCompile(`\bUSERNOTICE\b`)
	joinRegex *regexp.Regexp = regexp.MustCompile(`\bJOIN\b`)
	partRegex *regexp.Regexp = regexp.MustCompile(`\bPART\b`)
	cmdRegex  *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)
)

func EstablishWSConnection(channel string, username string, oath string) {
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
				chatMessage := MessageParser(rawIrcMessage)
				PrintChatMessage(chatMessage)
			} else if subRegex.MatchString(rawIrcMessage) {
				subMessage := SubParser(rawIrcMessage)
				PrintSubMessage(subMessage)
			} else if joinRegex.MatchString(rawIrcMessage) {
				// TODO: handle joins?
			} else if partRegex.MatchString(rawIrcMessage) {
				// TODO: handle parts?
			} else {
				fmt.Println(rawIrcMessage)
			}
		}
	}
}
