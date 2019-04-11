package main

import (
	"flag"
	"fmt"
	"syscall"
	"os"
	"os/signal"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	Perms = 251968
	ClientID = "237605108173635584"
	InvURL = "https://discordapp.com/oauth2/authorize?&client_id=%v&scope=bot&permissions=%v"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: cactusbot -t <token>")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}
	defer fmt.Println("\nGoodbye.")
	defer dg.Close() // close the session after Control-C

	fmt.Println("Cactusbot is now running. Press Control+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Client ready.")
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
