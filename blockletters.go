package main

import (
    "strings"
    "fmt"
)

func texttoemotes(t string) string {
    text := strings.ToLower(t)
    var output string
    for _, char := range(text) {
        if char >= 'a' && char <= 'z' {
            output += fmt.Sprintf(":regional_indicator_%c: ", char)
        } else if char == ' ' {
            output += "  "
        } else if char >= '0' && char <= '9' {
            output += fmt.Sprintf(":%v: ", numtext(char))
        } else {
            switch char {
                case '!':
                    output += ":exclamation: "
                case '?':
                    output += ":question: "
                case '#':
                    output += ":hash: "
                case '*':
                    output += ":asterisk: "
                case '$':
                    output += ":heavy_dollar_sign: "
                default:
                    output += string(char)
            }
        }
    }
    return output
}

func numtext(num rune) string {
    switch num {
        case '0':
            return "zero"
        case '1':
            return "one"
        case '2':
            return "two"
        case '3':
            return "three"
        case '4':
            return "four"
        case '5':
            return "five"
        case '6':
            return "six"
        case '7':
            return "seven"
        case '8':
            return "eight"
        case '9':
            return "nine"
        default:
            return string(num)
    }
}
