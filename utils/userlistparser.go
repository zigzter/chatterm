package utils

import (
	"strings"

	"github.com/zigzter/chatterm/types"
)

func UserListParser(input string) types.UserListMessage {
	var usernames []string
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if idx := strings.Index(line, ":"); idx != -1 {
			namesPart := line[idx+1:]
			namesList := strings.Split(namesPart, ":")
			// The user list message includes both 353 and 366 messages.
			// 353 lists users, 366 is an "End of /NAMES list" message,
			// so we only want to parse 353 messages.
			is353 := strings.Index(namesList[0], "353")
			if is353 != -1 {
				names := strings.Split(namesList[1], " ")
				usernames = append(usernames, names...)
			}
		}
	}
	return types.UserListMessage{Users: usernames}
}
