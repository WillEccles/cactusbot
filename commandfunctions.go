package main

import (
	"github.com/bwmarrin/discordgo"
	"regexp"
	"log"
	"math/rand"
	"fmt"
	"time"
	"strconv"
	"syscall"
)

func oodlehandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	re := regexp.MustCompile(`(?i)^c(actus)?\s+oodle\s+`)
	cleanmsg := re.ReplaceAllString(msg.Content, "")
	_, err := s.ChannelMessageSend(msg.ChannelID, oodle(cleanmsg))
	if err != nil {
		log.Printf("Error in oodlehandler:\n%v\n", err)
	}
}

func oodlettshandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	re := regexp.MustCompile(`(?i)^c(actus)?\s+oodletts\s+`)
	cleanmsg := re.ReplaceAllString(msg.Content, "")

	// check that the user has permission to use TTS, otherwise this will go poorly
	perms, err := s.State.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
	if err != nil {
		_, err = s.ChannelMessageSend(msg.ChannelID, "Something went wrong, please try again later. Sorry! :(")
	} else {
		if (perms & discordgo.PermissionSendTTSMessages) > 0 {
			_, err = s.ChannelMessageSendTTS(msg.ChannelID, oodle(cleanmsg))
		} else {
			_, err = s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Sorry <@%v>, you don't have permission to use TTS. Here's a normal one:\n%v", msg.Author.ID, oodle(cleanmsg)))

		}
	}
	
	if err != nil {
		log.Printf("Error in oodlettshandler:\n%v\n", err)
	}
}

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

func coinfliphandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	val := r1.Intn(2) // get a random number in [0, 2), so either 0 or 1
	var result string
	if val == 0 {
		result = "Heads!"
	} else {
		result = "Tails!"
	}
	_, err := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("**%v**", result))
	if err != nil {
		log.Printf("Error in coinfliphandler:\n%v\n", err)
	}
}

func blocklettershandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	re := regexp.MustCompile(`(?i)^c(actus)?\s+bl(ockletters)?\s+`)
	cleanmsg := re.ReplaceAllString(msg.Content, "")
	_, err := s.ChannelMessageSend(msg.ChannelID, texttoemotes(cleanmsg))
	if err != nil {
		log.Printf("Error in blocklettershandler:\n%v\n", err)
	}
}

func invitehandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	inv := fmt.Sprintf("Use this link to invite me to your server: " + InvURL, ClientID, Perms)
	_, err := s.ChannelMessageSend(msg.ChannelID, inv)
	if err != nil {
		log.Printf("Error in invitehandler:\n%v\n", err)
	}
}

func helphandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	embedcolor := s.State.UserColor(s.State.User.ID, msg.ChannelID)
	embed := HelpEmbed
	embed.Color = embedcolor

	_, err := s.ChannelMessageSendEmbed(msg.ChannelID, &embed)
	if err != nil {
		log.Printf("Error in helphandler:\n%v\n", err)
	}
}

func shutdownhandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	_, err := s.ChannelMessageSend(msg.ChannelID, "Goodbye!")
	if err != nil {
		log.Printf("Error in shutdownhandler:\n%v\n", err)
	}

	SigChan <- syscall.SIGINT
}

func srchandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	srcembed := &discordgo.MessageEmbed{
		URL: RepoURL,
		Color: s.State.UserColor(s.State.User.ID, msg.ChannelID),
		Title: "Github: willeccles/cactusbot",
		Description: "The source code for the cactus bot!",
		Author: &discordgo.MessageEmbedAuthor {
			URL: "https://eccles.dev",
			Name: "Will Eccles (a tiny cactus)",
			IconURL: "https://eccles.dev/imgs/avatar.jpg",
		},
		Fields: []*discordgo.MessageEmbedField {
			&discordgo.MessageEmbedField{
				Name: "Details",
				Value: "**Language:** go\n**Library:** discordgo",
			},
		},
	}
	_, err := s.ChannelMessageSendEmbed(msg.ChannelID, srcembed)
	if err != nil {
		log.Printf("Error in srchandler:\n%v\n", err)
	}
}

func xkcdhandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	re := regexp.MustCompile(`\D`) // just get rid of anything that's not a number
	cleanmsg := re.ReplaceAllString(msg.Content, "")
	comicnum := 0
	if cleanmsg != "" {
		comicnum, _ = strconv.Atoi(cleanmsg)
	}

	n, t, a, i, e := GetXkcd(comicnum)
	if e.ErrType != 0 {
		switch e.ErrType {
			case XkcdNotFound:
				_, err := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: xkcd #%v doesn't exist.", comicnum))
				if err != nil {
					log.Printf("Error in xkcdhandler:\n%v\n", err)
				}
			case XkcdNetworkErr:
				_, err := s.ChannelMessageSend(msg.ChannelID, "Error getting xkcd info. Please try again later.")
				if err != nil {
					log.Printf("Error in xkcdhandler:\n%v\n", err)
				}
			case XkcdOtherErr:
				_, err := s.ChannelMessageSend(msg.ChannelID, "Error getting xkcd info. Please try again later.")
				if err != nil {
					log.Printf("Error in xkcdhandler:\n%v\n", err)
				}
		}
		return
	}

	url := "https://xkcd.com/"
	if comicnum != 0 {
		url += string(comicnum)
	}

	xkcdembed := &discordgo.MessageEmbed{
		URL: url,
		Color: s.State.UserColor(s.State.User.ID, msg.ChannelID),
		Title: fmt.Sprintf("#%v: **%v**", n, t),
		Image: &discordgo.MessageEmbedImage{
			URL: i,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: a,
		},
	}

	_, err := s.ChannelMessageSendEmbed(msg.ChannelID, xkcdembed)
	if err != nil {
		log.Printf("Error in xkcdhandler:\n%v\n", err)
	}
}
