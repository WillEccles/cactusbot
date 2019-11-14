package main

import (
	"regexp"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"strings"
)

func (carg CommandArg) String() string {
	if carg.Required {
		return fmt.Sprintf("<%v>", carg.Title)
	} else {
		return fmt.Sprintf("[%v]", carg.Title)
	}
}

func (cmd *Command) Handle(msg *discordgo.MessageCreate, s *discordgo.Session) {
	if cmd.AdminOnly && (!Config.IsAdmin(msg.Author.ID) && msg.Author.ID != Config.ControllerID ) {
		return
	}

	cmd.Handler(msg, s)
}

var CommandCategories = map[string]*struct{
	Title string
	Cmds []*Command
}{
	"text": {
		Title: ":page_facing_up: Text",
	},
	"fun": {
		Title: ":game_die: Fun",
	},
	"util": {
		Title: ":wrench: Utility",
	},
}

// go iterates over maps (using range()) in a random order, so this is used to combat that
var CmdCatOrder = []string{ "fun", "text", "util" }

var Commands = []Command {
	{
		Name: "oodle",
		Args: []CommandArg {
			{
				Title: "message",
				Required: true,
			},
		},
		Description: "Replaces every vowel in `message` with 'oodle' or 'OODLE', depending on whether or not it's a capital.",
		Examples: []string{
			"`c oodle I am a bot.` returns \"OODLE oodlem oodle boodlet.\"",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+oodle\s+.*[aeiou].*`),
		Category: "text",
		Handler: oodlehandler,
	},
	{
		Name: "oodletts",
		Args: []CommandArg {
			{
				Title: "message",
				Required: true,
			},
		},
		Description: "Works the same as `oodle`, but responds with a TTS message. Requires the user to have permission to use TTS.",
		Examples: []string{
			"`c oodletts I am a bot.` returns \"OODLE oodlem oodle boodlet.\"",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+oodletts\s+.*[aeiou].*`),
		Category: "text",
		Handler: oodlettshandler,
	},
	{
		Name: "coinflip",
		Description: "Flips a coin.",
		Examples: []string{
			"`c coinflip` returns either Heads or Tails.",
		},
		Aliases: []string {
			"cf",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+(coinflip|cf)`),
		Category: "fun",
		Handler: coinfliphandler,
	},
	{
		Name: "roll",
		Args: []CommandArg {
			{
				Title: "count",
				Required: false,
			},
			{
				Title: "sides",
				Required: false,
			},
		},
		Description: "Rolls some dice. If no arguments are supplied, a standard 6-sided die is rolled. If there's only one argument, it's `sides`. If there are two, the number of dice comes first, then the number of sides. These can also be separated by a 'd', like `2d10`, and if only sides is specified, it can be `d10`.",
		Examples: []string{
			"`c roll` returns 1-6",
			"`c roll 20` returns 1-20",
			"`c roll d10` returns 1-10",
			"`c roll 2d20` returns 1-20 for two dice",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+roll(\s+(\d+)?(\s*d\s*|\s+)?\d+)?`),
		Category: "fun",
		Handler: rollhandler,
	},
	{
		Name: "blockletters",
		Args: []CommandArg {
			{
				Title: "message",
				Required: true,
			},
		},
		Description: "Converts as much of `message` as possible into block letters using emoji.",
		Examples: []string{
			"`c blockletters Something` returns \"Something\" written in blockletters.",
		},
		Aliases: []string {
			"bl",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+bl(ockletters)?\s+\S+`),
		Category: "text",
		Handler: blocklettershandler,
	},
	{
		Name: "xkcd",
		Args: []CommandArg {
			{
				Title: "number",
				Required: false,
			},
		},
		Description: "Gets either the most recent xkcd or the xkcd with the given `number`.",
		Examples: []string{
			"`c xkcd` embeds the most recent xkcd.",
			"`c xkcd 327` embeds the Little Bobby Tables xkcd.",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+xkcd(\s+\d+)?`),
		Category: "fun",
		Handler: xkcdhandler,
	},
	{
		Name: "invite",
		Description: "Creates a discord invite link to add this bot to another server.",
		Aliases: []string {
			"inv",
		},
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+inv(ite)?`),
		Category: "util",
		Handler: invitehandler,
	},
	{
		Name: "help",
		Description: "Displays this help message.",
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+help.*`),
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
		Category: "util",
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
	{
		Name: "sowner",
		Description: "Returns the user who owns the server.",
		AdminOnly: true,
		Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+sowner`),
		Handler: sownerhandler,
	},
    {
        Name: "echo",
        Description: "Echos back what you say.",
        AdminOnly: true,
        Pattern: regexp.MustCompile(`(?i)^c(actus)?\s+echo\s+\S+`),
        Handler: echohandler,
    },
    {
        Name: "lol",
        Description: "Looks up league-related things.",
        AdminOnly: true, // TODO
        Pattern: regexp.MustCompile(`(?i)^c\s+lol\s+(profile|status)\s+\S+`),
        Handler: lolhandler,
    },
}

func InitHelpEmbed(embed *discordgo.MessageEmbed) {
	for i, cmd := range(Commands) {
		if cmd.Category == "" {
			continue
		}
		CommandCategories[cmd.Category].Cmds = append(CommandCategories[cmd.Category].Cmds, &(Commands[i]))
	}

	// prepare a help embed to reduce CPU load later on
	embed.Title = "Command List"
	embed.Description = "You should begin each command with `cactus` or simply `c`.\nFor example: `cactus help` or `c help`.\nFor info about a particular command, use `c help [command]`."

	for _, catname := range(CmdCatOrder) {
		cat := CommandCategories[catname]

		newfield := discordgo.MessageEmbedField{
			Name: cat.Title,
			Inline: false,
		}

		for _, cmd := range(cat.Cmds) {
			// only show non-admin commands
			if cmd.AdminOnly {
				continue
			}

			newfield.Value += fmt.Sprintf("`%v` ", cmd.Name)
		}

		embed.Fields = append(embed.Fields, &newfield)
	}
}

func InitCommandEmbeds(m map[string]*discordgo.MessageEmbed) {
	for _, cmd := range(Commands) {
		m[cmd.Name] = &discordgo.MessageEmbed{}
		m[cmd.Name].Title = "`" + cmd.Name
		for _, arg := range(cmd.Args) {
			m[cmd.Name].Title += fmt.Sprintf(" %s", arg)
		}
		m[cmd.Name].Title += "`"

		if cmd.Args != nil {
			m[cmd.Name].Footer = &discordgo.MessageEmbedFooter{
				Text: "Arguments in <angle brackets> are required, arguments in [square brackets] are optional.",
			}
		}

		m[cmd.Name].Description = cmd.Description

		if cmd.Examples != nil {
			m[cmd.Name].Fields = append(m[cmd.Name].Fields, &discordgo.MessageEmbedField{
				Name: "Examples",
				Value: strings.Join(cmd.Examples, "\n"),
				Inline: false,
			})
		}

		if cmd.Aliases != nil {
			m[cmd.Name].Fields = append(m[cmd.Name].Fields, &discordgo.MessageEmbedField{
				Name: "Aliases",
				Value: fmt.Sprintf("`%v`", strings.Join(cmd.Aliases, "` `")),
				Inline: false,
			})

			for _, a := range(cmd.Aliases) {
				m[a] = m[cmd.Name]
			}
		}
	}
}
