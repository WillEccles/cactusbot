package main

import (
	"fmt"
	"syscall"
	"os"
	"os/signal"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	Perms = 251968
	InvURL = "<https://discordapp.com/oauth2/authorize?&client_id=%v&scope=bot&permissions=%v>"
	RepoURL = "https://github.com/willeccles/cactusbot"
)

var HelpEmbed discordgo.MessageEmbed
var SigChan chan os.Signal
var Config Configuration

var CommandEmbeds map[string]*discordgo.MessageEmbed


func init() {
	log.SetPrefix("[Cactusbot] ")
	log.Println("init: loading config")
	Config = LoadConfig()
	log.Println("init: creating help embeds")
	InitHelpEmbed(&HelpEmbed)
	CommandEmbeds = make(map[string]*discordgo.MessageEmbed)
	InitCommandEmbeds(CommandEmbeds)
}

func main() {
	if Config.DiscordToken == "" {
		log.Println("Please provide a discord token in your config.json file.")
		return
	}
	if Config.DiscordClientID == "" {
		log.Println("Please provide a discord client ID in your config.json file.")
		return
	}
	if Config.DebugChannel == "" {
		log.Println("Please provide a debug channel ID in your config.json file.")
		return
	}
	if len(Config.AdminIDs) == 0 {
		log.Println("Please provide at least one admin ID in your config.json file.")
		return
	}
	if Config.ControllerID == "" {
		log.Println("Controller ID not found in config.json, assuming no controller.")
	}

	dg, err := discordgo.New("Bot " + Config.DiscordToken)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(connect)
	dg.AddHandler(resume)
	dg.AddHandler(disconnect)

	err = dg.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
		return
	}
	defer fmt.Println("\nGoodbye.")
	defer dg.Close() // close the session after Control-C

	SigChan = make(chan os.Signal)
	signal.Notify(SigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-SigChan
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Client ready.")

	// set the status to "watching you"
	i := 0
	usd := discordgo.UpdateStatusData{
		IdleSince: &i,
		AFK: false,
		Status: "online",
		Game: &discordgo.Game {
			Name: "you OwO",
			Type: discordgo.GameTypeWatching,
		},
	}

	err := s.UpdateStatusComplex(usd)
	if err != nil {
		log.Printf("Error in ready:\n%v\n", err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore the bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, cmd := range(Commands) {
		if cmd.Pattern.MatchString(m.Content) {
			cmd.Handle(m, s)
			break
		}
	}
}

func connect(s *discordgo.Session, event *discordgo.Connect) {
	log.Println("Client connected.")
}

func disconnect(s *discordgo.Session, event *discordgo.Disconnect) {
	log.Println("Client disconnected!")
}

func resume(s *discordgo.Session, event *discordgo.Resumed) {
	log.Println("Resumed, attempting to send debug message.")
	_, err := s.ChannelMessageSend(Config.DebugChannel, fmt.Sprintf("Just recovered from error(s)!\n```\n%v\n```", strings.Join(event.Trace, "\n")))
	if err != nil {
		log.Printf("Error in resume (this is awkward):\n%v\n", err)
	}
}
