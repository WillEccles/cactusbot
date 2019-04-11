package main

import (
	"regexp"
	"github.com/bwmarrin/discordgo"
)

type MsgHandler func(*discordgo.MessageCreate, *discordgo.Session)

type Command struct {
	Pattern *regexp.Regexp
	Name string
	Description string
	Aliases []string
	Handler MsgHandler
	AdminOnly bool
}

func (cmd *Command) Handle(msg *discordgo.MessageCreate, s *discordgo.Session) {
	if cmd.AdminOnly && msg.Author.ID != AdminID {
		return
	}

	cmd.Handler(msg, s)
}

var Commands = []Command {
	{
		Name: "oodle",
		Description: "Replaces every vowel in a message with 'oodle' or 'OODLE', depending on whether or not it's a capital.",
		Pattern: regexp.MustCompile(`^(?i)c(actus)?\s+oodle\s+.*[aeiou].*`),
		Handler: oodlehandler,
	},
	{
		Name: "oodletts",
		Description: "Works the same as `oodle`, but responds with a TTS message.",
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+oodletts\s+.*[aeiou].*`),
		Handler: oodlettshandler,
	},
	{
		Name: "coinflip",
		Description: "Flips a coin.",
		Aliases: []string {
			"cf",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+(coinflip|cf)`),
		Handler: coinfliphandler,
	},
	{
		Name: "blockletters",
		Description: "Converts as much of a message as possible into block letters using emoji.",
		Aliases: []string {
			"bl",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+bl(ockletters)?\s+\S+`),
		Handler: blocklettershandler,
	},
	{
		Name: "xkcd",
		Description: "Displays either today's xkcd or the specified xkcd. For today's, simply use `xkcd`. For a specific one, use `xkcd <number>`.",
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+xkcd(\s+\d+)?`),
		Handler: xkcdhandler,
	},
	{
		Name: "invite",
		Description: "Creates a discord invite link for to add this bot to another server.",
		Aliases: []string {
			"inv",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+inv(ite)?`),
		Handler: invitehandler,
	},
	{
		Name: "help",
		Description: "Displays this help message.",
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+help`),
		Handler: helphandler,
	},
	{
		Name: "source",
		Description: "Gives you a link to my source code.",
		Aliases: []string {
			"src",
			"git",
			"repo",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+(source|src|git|repo)`),
		Handler: srchandler,
	},
	{
		Name: "shutdown",
		Description: "Shuts down the bot.",
		Aliases: []string {
			"sd",
		},
		AdminOnly: true,
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+(shutdown|sd)`),
		Handler: shutdownhandler,
	},
}
