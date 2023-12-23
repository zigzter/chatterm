package utils

import (
    "fmt"
)

func PrintMessage(message ChatMessage) {
    icon := ""
    bgColor := ""
    if message.IsFirstMessage {
        bgColor = "\033[41m"
    }
    modIcon := "\033[32m󰓥"
    vipIcon := "\033[35m󰮊"
    if message.IsMod {
        icon = modIcon
    } else if message.IsVIP {
        icon = vipIcon
    }

    color := ParseHexColor(message.Color)
    resetCode := "\033[0m"
    defaultTextColor := "\033[39m"
    fmt.Printf(
        "%s%s\033[38;2;%d;%d;%dm%s%s%s: %s%s\n", // Format string
        bgColor,                                // Background color
        icon,                                   // Icon (mod/vip)
        color.R, color.G, color.B,              // Text color (RGB) for DisplayName
        message.DisplayName,                    // Display name
        defaultTextColor,                       // Set text color to default before Message
        resetCode,                              // Reset text color after DisplayName
        message.Message,                        // Message
        resetCode,                              // Reset formatting at the end
    )
}

