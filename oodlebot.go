package main

import (
	"flag"
	"fmt"
	"syscall"
	"os"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: oodlebotgo -t <token>")
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
	defer dg.Close() // close the session after Control-C

	fmt.Println("Oodlebot is now running. Press Control+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// TODO
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO
}
