package utils

import (
    "fmt"
)

func PrintMessage(message ChatMessage) {
    icon := ""
    modIcon := "󰓥 "
    vipIcon := " "
    if message.IsMod {
        icon = modIcon
    } else if message.IsVIP {
        icon = vipIcon
    }
    color := ParseHexColor(message.Color)
    // fmt.Printf("%+v\n", message)
    // fmt.Printf("%s %s: %s \n", icon, message.DisplayName, message.Message)
    fmt.Printf("%s\033[38;2;%d;%d;%dm%s\033[0m: %s\n", icon, color.R, color.G, color.B, message.DisplayName, message.Message)
}

