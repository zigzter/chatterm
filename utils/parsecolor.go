package utils

import (
    "image/color"
)

func ParseHexColor(s string) (c color.RGBA) {
    c.A = 0xff
    defaultColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

    if len(s) == 0 || s[0] != '#' {
        return defaultColor
    }

    hexToByte := func(b byte) byte {
        switch {
        case b >= '0' && b <= '9':
            return b - '0'
        case b >= 'a' && b <= 'f':
            return b - 'a' + 10
        case b >= 'A' && b <= 'F':
            return b - 'A' + 10
        }
        return 0
    }

    switch len(s) {
    case 7:
        c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
        c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
        c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
    case 4:
        c.R = hexToByte(s[1]) * 17
        c.G = hexToByte(s[2]) * 17
        c.B = hexToByte(s[3]) * 17
    default:
        return defaultColor
    }
    return
}

