package main

import (
    "regexp"
    "github.com/bwmarrin/discordgo"
)

/* Commands */

type MsgHandler func(*discordgo.MessageCreate, *discordgo.Session)

type Command struct {
    Pattern     *regexp.Regexp
    Name        string
    Args        []CommandArg
    Examples    []string
    Description string
    Aliases     []string
    Handler     MsgHandler
    Category    string // if "" the command won't be listed in help menu
    AdminOnly   bool
}

type CommandArg struct {
    Title       string
    Required    bool
}
