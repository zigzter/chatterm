# Chatterm

Twitch chat in the terminal. Currently just spits out new chats as they come in. Goal is to be able to chat as well as perform moderator actions.

Current flow (requires Go if cloning):

1. Either clone and `go build .`, or download the binary under releases
2. Grab a token from here: https://twitchapps.com/tmi/
3. Run `./chatterm config` and enter your username and Oauth token when prompted
3. `./chatterm connect -c channelName`

Todo:

- [x] Connect to Twitch channel via flags
- [ ] Add ability to send chats
- [x] Store Oauth and username in a local config
- [ ] Customize chat output (show/hide badges, colors for first chatters, etc)
- [ ] Perform moderator actions

