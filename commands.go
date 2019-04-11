package main

import (
	"regexp"
	"github.com/bwmarrin/discordgo"
)

type MsgHandler func(*discordgo.Message, *discordgo.Session)

type Command struct {
	Pattern *regexp.Regexp
	Name string
	Description string
	Aliases []string
	Handler MsgHandler
}

var Commands = []Command {
	{
		Name: "oodle",
		Description: "Replaces every vowel in a message with 'oodle' or 'OODLE', depending on whether or not it's a capital.",
		Pattern: regexp.MustCompile(`^c(actus)? oodle\s+.*[aeiouAEIOU].*`),
		Handler: oodlehandler,
	},
	{
		Name: "oodletts",
		Description: "Works the same as ?oodle, but responds with a TTS message.",
		Pattern: regexp.MustCompile(`^c(actus)? oodletts\s+.*[aeiouAEIOU].*`),
		Handler: oodlettshandler,
	},
	{
		Name: "coinflip",
		Description: "Flips a coin.",
		Aliases: []string {
			"cf",
		},
		Pattern: regexp.MustCompile(`^c(actus)? coinflip`),
		Handler: coinfliphandler,
	},
	{
		Name: "blockletters",
		Description: "Converts as much of a message as possible into block letters using emoji.",
		Aliases: []string {
			"bl",
		},
		Pattern: regexp.MustCompile(`^c(actus)? bl(ockletters)?\s+\S+`),
		Handler: blocklettershandler,
	},
	{
		Name: "invite",
		Description: "Creates a discord invite link for to add this bot to another server.",
		Aliases: []string {
			"inv",
		},
		Pattern: regexp.MustCompile(`^c(actus)? inv(ite)?`),
		Handler: invitehandler,
	},
	{
		Name: "help",
		Description: "Displays the help message.",
		Pattern: regexp.MustCompile(`^c(actus)? help`),
		Handler: helphandler,
	},
}
