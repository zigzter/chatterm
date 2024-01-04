# Chatterm

Twitch chat in the terminal. Using Bubble Tea for the terminal UI, Gorilla WebSocket for the connection, & Viper for the config.

Current flow (requires Go if cloning):

1. Either clone and `go build .`, or download the binary under releases
2. Run `./chatterm`, follow prompt to open auth and input username and start auth process.
3. Once submitted, it should bring you back to the channel input view. Enter a channel to join and press enter.

Todo:

- [x] Connect to Twitch channel via flags
- [x] Add ability to send chats
- [x] Store Oauth and username in a local config
- [ ] Customize chat output (show/hide badges, colors for first chatters, etc)
- [ ] Allow saving chat to file
- [ ] Allow searching chat
- [ ] Perform moderator actions

Twitch [announced](https://discuss.dev.twitch.com/t/deprecation-of-chat-commands-through-irc/40486) awhile back that mod commands via IRC will no longer work, and instead have to go through their API. This has put a bit of a damper on my moderator actions plan. I'll have to figure out a way to utilize the API while still letting the app be useable by other people.
