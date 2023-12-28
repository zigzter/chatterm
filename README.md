# Chatterm

Twitch chat in the terminal. Using Bubble Tea for the terminal UI, Gorilla WebSocket for the connection, & Viper for the config.

Current flow (requires Go if cloning):

1. Either clone and `go build .`, or download the binary under releases
2. Grab a token from here: https://twitchapps.com/tmi/
3. Run `./chatterm config` and enter your username and Oauth token when prompted (currently does not work, store chatterm.json manually in `~/.config/chatterm.json`)
4. `./chatterm connect -c channelName`

Todo:

- [x] Connect to Twitch channel via flags
- [x] Add ability to send chats
- [x] Store Oauth and username in a local config
- [ ] Customize chat output (show/hide badges, colors for first chatters, etc)
- [ ] Perform moderator actions

