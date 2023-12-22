package utils

import (
    "fmt"
)

func PrintMessage(message ChatMessage) {
    icon := ""
    modIcon := "󰓥"
    vipIcon := ""
    if message.IsMod {
        icon = modIcon
    } else if message.IsVIP {
        icon = vipIcon
    }
    // fmt.Printf("%+v\n", message)
    fmt.Printf("%s %s: %s \n", icon, message.DisplayName, message.Message)
}

