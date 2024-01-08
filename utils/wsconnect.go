package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/types"
)

var (
	msgRegex             *regexp.Regexp = regexp.MustCompile(`\bPRIVMSG\b`)
	subRegex             *regexp.Regexp = regexp.MustCompile(`\bUSERNOTICE\b`)
	globalUserStateRegex *regexp.Regexp = regexp.MustCompile(`\bGLOBALUSERSTATE\b`)
	joinRegex            *regexp.Regexp = regexp.MustCompile(`\bJOIN\b`)
	partRegex            *regexp.Regexp = regexp.MustCompile(`\bPART\b`)
	pingRegex            *regexp.Regexp = regexp.MustCompile(`\bPING\b`)
	cmdRegex             *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)
)

type WebSocketClient struct {
	Conn *websocket.Conn
}

func NewWebSocketClient() (*WebSocketClient, error) {
	socketUrl := "ws://irc-ws.chat.twitch.tv:80"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocketClient{Conn: conn}, nil
}

func (client *WebSocketClient) SendMessage(message []byte) error {
	return client.Conn.WriteMessage(websocket.TextMessage, message)
}

func EstablishWSConnection(client *WebSocketClient, channel string, username string, oath string, msgChan chan<- types.ChatMessageWrap) {
	err1 := client.SendMessage([]byte("CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands"))
	if err1 != nil {
		log.Println("Error writing message:", err1)
	}
	setPass := fmt.Sprintf("PASS %s", oath)
	client.SendMessage([]byte(setPass))
	setUser := fmt.Sprintf("NICK %s", username)
	client.SendMessage([]byte(setUser))
	joinChannel := fmt.Sprintf("JOIN #%s", channel)
	client.SendMessage([]byte(joinChannel))
	for {
		// Read messages and handle them
		messageType, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			return
		}
		if messageType == websocket.TextMessage {
			rawIrcMessage := strings.TrimSpace(string(message))
			switch {
			case msgRegex.MatchString(rawIrcMessage):
				chatMessage := MessageParser(rawIrcMessage)
				// TODO: find a better place to map the id
				sql := db.OpenDB()
				db.InsertUserMap(sql, chatMessage.DisplayName, chatMessage.UserId)
				msgChan <- types.ChatMessageWrap{ChatMsg: chatMessage}
			case subRegex.MatchString(rawIrcMessage):
				subMessage := SubParser(rawIrcMessage)
				msgChan <- types.ChatMessageWrap{ChatMsg: subMessage}
			case joinRegex.MatchString(rawIrcMessage):
			case partRegex.MatchString(rawIrcMessage):
			case pingRegex.MatchString(rawIrcMessage):
				client.SendMessage([]byte("PONG :tmi.twitch.tv"))
			case globalUserStateRegex.MatchString(rawIrcMessage):
				StoreUserState(rawIrcMessage)
			default:
				log.Println(rawIrcMessage)
			}
		}
	}
}
