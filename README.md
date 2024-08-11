# Chatterm

Twitch chat in the terminal, with moderator actions.

![Chat app preview image](./chat_view.png)

![Chat app preview image](./info_view.png)

![Chat app channel input image](./channel_input.png)

## ⚡ Features
- View and send messages.
- Built-in authentication via Twitch.
- Ban & timeout users, clear chat, query user info
- Search received messsages, powered by the FTS5 extension of SQLite
- Username autocomplete on `@` mentions and `/` commands, triggered by pressing tab.

## 🔧 Tech Used
- Bubble Tea for the terminal UI
- Gorilla WebSocket for the websocket connection
- SQLite for storage
- Viper for configuration

## 🖥️ Supported Platforms
Currently only tested on Linux. macOS might work if you build it yourself.

## 📦 Installation
Cloning option (requires Go):
1. `git clone https://github.com/zigzter/chatterm.git`
2. `cd chatterm`
3. `go build --tags "fts5" .`
4. `./chatterm`

If you have to re-auth for whatever reason, and the oauth request gets stuck loading, try removing the connection from Twitch Connections, then retry the auth.

Downloading binary option:
Simply download the binary and run `./chatterm`

## 🚀 Supported Commands
- Ban a user: `/ban username`
- Timeout a user: `/ban username timeInSeconds`
- Search for messages: `/search something to search`
    - `from:username` restricts search to messages sent by that user
    - `channel:channelName` restricts search to messages in that channel
    - Supports wildcards and is case insensitive: `kek*` will find `KEKW`
    - Supports `AND`, `NOT`, and `OR` keywords. `/search tf2 OR ow2` will find all messages that contain at least one of those
    - This search queries the local DB, so it will only find messages that you've received while in the chat rooms
- Clear chat: `/clear` (this is the moderator clear, not a local one)
- Get a user's info: `/info username`
- Send an announcement: `/announcement something to announce`
- Give a streamer a shoutout: `/shoutout username` (untested, only works if the channel you're in is live)
- Warn a user, requiring them to acknowledge the warning before they can resume chatting: `/warn username reason`
- Delete all locally stored chat messages: `/clearall`
- Enable/disable shield mode: `/shield on|off`
- Watch a user (highlight their messages): `/watch username`
    - To remove the user, re-run the command
    - You can also manually edit the `$HOME/.config/chatterm.json` file to add/remove users under the `watched-users` key

