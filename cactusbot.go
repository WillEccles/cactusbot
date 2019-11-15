package main

import (
	"fmt"
	"syscall"
	"os"
	"os/signal"
	"log"
	"regexp"
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

var LeagueData LeagueHelper
var EnableLOL = true

var DadMatcher = regexp.MustCompile(`(?i)^i(['’]?m|\s+am)\s+\S`)
var DadReplacer = regexp.MustCompile(`(?i)^i(['’]?m|\s+am)\s+`)
var DadSanitizer = regexp.MustCompile(`(?i)@+(everyone|here)`)
var DadEnabler = regexp.MustCompile(`(?i)^c\s+dad\s+(on|off)`)
var EnableDad = false

func init() {
	log.SetPrefix("[Cactusbot] ")
	log.Println("init: loading config")
	Config = LoadConfig()
	//go WriteConfig(Config) // if config is out of date, this updates it
    if Config.LeagueToken == "" {
        log.Println("League token not found; 'lol' commands will be disabled.")
        EnableLOL = false
    } else {
        if !LeagueData.Init(Config.LeagueToken) {
            // TODO disable league commands
            log.Println("Error initializing league data; league commands will be disabled")
            EnableLOL = false
        }
    }
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
		log.Println("ControllerID not found in config.json, assuming no controller.")
	}
    if Config.LogChannel == "" || Config.LogWebhookID == "" || Config.LogWebhookToken == "" {
        log.Println("Channel logging not enabled.")
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

    if EnableLOL {
        go LeagueData.UpdateRoutine()
    }

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
	// ignore bots' messages
	if m.Author.Bot {
		return
	}

    if Config.IsAdmin(m.Author.ID) && DadEnabler.MatchString(m.Content) {
        parts := strings.Split(m.Content, " ")
        parts[2] = strings.ToLower(parts[2])
        if parts[2] == "off" {
            EnableDad = false
        } else if parts[2] == "on" {
            EnableDad = true
        }
        _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```golang\nEnableDad = %t```", EnableDad))
        if err != nil {
            log.Printf("Error in messageCreate:\n%v\n", err)
        }
    }

    if EnableDad && DadMatcher.MatchString(m.Content) {
        response := DadReplacer.ReplaceAllString(m.Content, "")
        response = DadSanitizer.ReplaceAllString(response, "$1")
        _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi %s, I'm cactusbot!", response))
        if err != nil {
            log.Printf("Error in messageCreate:\n%v\n", err)
        }
    }

	for _, cmd := range(Commands) {
		if cmd.Pattern.MatchString(m.Content) {
			cmd.Handle(m, s)
			break
		}
	}

    if !(Config.LogChannel == "" || Config.LogWebhookID == "" || Config.LogWebhookToken == "") {
        if m.ChannelID == Config.LogChannel {
            whp := discordgo.WebhookParams{
                Content: strings.ReplaceAll(m.Content, "@", ""),
                Username: m.Author.Username,
                AvatarURL: m.Author.AvatarURL(""),
            }
            if len(m.Attachments) > 0 {
                for _, a := range m.Attachments {
                    whp.Embeds = append(whp.Embeds, &discordgo.MessageEmbed{
                        Image: &discordgo.MessageEmbedImage{
                            URL: a.URL,
                        },
                    })
                }
            }
            err := s.WebhookExecute(Config.LogWebhookID, Config.LogWebhookToken, false, &whp)
            if err != nil {
                log.Println("Error in messageCreate:\n%v\n", err)
            }
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

	err = s.UpdateStatusComplex(usd)
	if err != nil {
		log.Printf("Error in resume (this is awkward):\n%v\n", err)
	}
}
