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
	roomStateRegex       *regexp.Regexp = regexp.MustCompile(`\bROOMSTATE\b`)
	joinRegex            *regexp.Regexp = regexp.MustCompile(`\bJOIN\b`)
	partRegex            *regexp.Regexp = regexp.MustCompile(`\bPART\b`)
	pingRegex            *regexp.Regexp = regexp.MustCompile(`\bPING\b`)
	listUsersRegex       *regexp.Regexp = regexp.MustCompile(`\b353\b`)
	cmdRegex             *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)
)

type WebSocketClient struct {
	Conn *websocket.Conn
}

// NewWebSocketClient creates and returns the websocket client
func NewWebSocketClient() (*WebSocketClient, error) {
	socketUrl := "ws://irc-ws.chat.twitch.tv:80"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocketClient{Conn: conn}, nil
}

// SendMessage is a wrapper for sending IRC messages to the IRC server
func (client *WebSocketClient) SendMessage(message []byte) error {
	return client.Conn.WriteMessage(websocket.TextMessage, message)
}

// EstablishWSConnection sends the authentication information to Twitch,
// then listens for and processes incoming IRC messages
func EstablishWSConnection(client *WebSocketClient, channel string, username string, oath string, msgChan chan<- types.ParsedIRCMessage) {
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
				msgChan <- types.ParsedIRCMessage{Msg: chatMessage}
			case subRegex.MatchString(rawIrcMessage):
				subMessage := SubParser(rawIrcMessage)
				msgChan <- types.ParsedIRCMessage{Msg: subMessage}
			case joinRegex.MatchString(rawIrcMessage):
			case partRegex.MatchString(rawIrcMessage):
			case listUsersRegex.MatchString(rawIrcMessage):
				userListMessage := UserListParser(rawIrcMessage)
				msgChan <- types.ParsedIRCMessage{Msg: userListMessage}
			case pingRegex.MatchString(rawIrcMessage):
				client.SendMessage([]byte("PONG :tmi.twitch.tv"))
			case globalUserStateRegex.MatchString(rawIrcMessage):
				StoreUserState(rawIrcMessage)
			case roomStateRegex.MatchString(rawIrcMessage):
				StoreRoomState(rawIrcMessage)
			default:
				log.Println(rawIrcMessage)
			}
		}
	}
}
