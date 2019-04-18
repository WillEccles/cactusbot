package main

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"golang.org/x/text/language"
    "golang.org/x/text/language/display"
)

const (
	StreamStatusURL = "https://api.twitch.tv/kraken/streams/%s?client_id=%s"
	ChannelInfoURL = "https://api.twitch.tv/kraken/channels/%s?client_id=%s"
	ChannelFollowURL = "https://api.twitch.tv/kraken/users/%s/follows/channels/%s?client_id=%s"
)

var TwitchErrorEmbed = &discordgo.MessageEmbed{
	Title: "Error",
	Description: "Error getting data from Twitch. Please try again later.",
	Color: 0xCE0000,
}
var TwitchNoSuchChannelEmbed = &discordgo.MessageEmbed{
	Title: "Error",
	Description: "That channel does not exist.",
	Color: 0xCE0000,
}
var TwitchFollowAgeErrorEmbed = &discordgo.MessageEmbed{
	Title: "Error",
	Description: "Either one of the users doesn't exist or that user doesn't follow that channel.",
	Color: 0xCE0000,
}

func GetStreamStatusEmbed(stream string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}
	
	resp, err := http.Get(fmt.Sprintf(StreamStatusURL, stream, Config.TwitchClientID))
	if err != nil {
		if resp.StatusCode == 404 {
			return TwitchNoSuchChannelEmbed
		} else {
			log.Printf("GetStreamStatusEmbed: %v\n", err)
			return TwitchErrorEmbed
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GetStreamStatusEmbed: %v\n", err)
		return TwitchErrorEmbed
	}
	
	var status StreamStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Printf("GetStreamStatusEmbed: %v\n", err)
		return TwitchErrorEmbed
	}

	if status.Stream == nil {
		return nil // returning nil indicates to the caller that they should instead try using GetChannelInfoEmbed
	}

	tag := language.Make(status.Stream.Channel.Language)
	langstr := fmt.Sprintf("%s (%s)",
		display.English.Tags().Name(tag),
		display.Self.Name(tag))

	embed.Title = "Details"
	embed.Description = "**Status:** " + strings.Title(status.Stream.StreamType)
	embed.Description += "\n**Title:** " + strings.TrimSuffix(status.Stream.Channel.Status, "\n")
	embed.Description += "\n**Game:** " + status.Stream.Game
	embed.Description += "\n**Viewers:** " + strconv.Itoa(status.Stream.Viewers)
	embed.Description += "\n**Uptime:** " + GetUptimeString(status.Stream)
	embed.Description += "\n**Language:** " + langstr
	embed.Color = 0x6441A5 // twitch purple
	embed.Author = &discordgo.MessageEmbedAuthor{
		URL: status.Stream.Channel.URL,
		Name: status.Stream.Channel.DisplayName,
		IconURL: status.Stream.Channel.Avatar,
	}
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: status.Stream.Previews.Large,
		Width: 80,
		Height: 45,
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "Statistics",
		Value: fmt.Sprintf("**Followers:** %v\n**Views:** %v\n**Account Created:** %v\n**Partner:** %t\n**Mature:** %t", status.Stream.Channel.Followers, status.Stream.Channel.Views, GetChannelAgeString(status.Stream.Channel), status.Stream.Channel.Partner, status.Stream.Channel.Mature),
		Inline: false,
	})
	
	return embed
}

func GetChannelInfoEmbed(channel string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}
	
	resp, err := http.Get(fmt.Sprintf(ChannelInfoURL, channel, Config.TwitchClientID))
	if err != nil {
		if resp.StatusCode == 404 {
			return TwitchNoSuchChannelEmbed
		} else {
			log.Printf("GetChannelInfoEmbed: %v\n", err)
			return TwitchErrorEmbed
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GetChannelInfoEmbed: %v\n", err)
		return TwitchErrorEmbed
	}
	
	var res TwitchChannel
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("GetChannelInfoEmbed: %v\n", err)
		return TwitchErrorEmbed
	}

	if res.Error != "" {
		return TwitchNoSuchChannelEmbed // returning nil indicates to the caller that they should instead try using GetChannelInfoEmbed
	}

	embed.Title = "Info"
	embed.Color = 0x6441A5
	embed.Description = fmt.Sprintf("**Last Game:** %s\n**Followers:** %v\n**Views:** %v\n**Account Created:** %v\n**Partner:** %t\n**Mature:** %t", res.Game, res.Followers, res.Views, GetChannelAgeString(&res), res.Partner, res.Mature)

	embed.Author = &discordgo.MessageEmbedAuthor{
		URL: "https://twitch.tv/" + channel,
		Name: res.DisplayName,
		IconURL: res.Avatar,
	}

	return embed
}

func GetUptimeString(stream *TwitchStream) string {
	created, _ := time.Parse("2006-01-02T15:04:05Z", stream.CreatedAt)
	dur := time.Since(created).Round(time.Second)
	return dur.String()
}

func GetChannelAgeString(channel *TwitchChannel) string {
	created, _ := time.Parse("2006-01-02T15:04:05Z", channel.CreatedAt)
	return created.Format("Jan _2 2006 15:04:05") + " UTC"
}

func GetFollowAgeEmbed(user string, channel string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}
	
	resp, err := http.Get(fmt.Sprintf(ChannelFollowURL, user, channel, Config.TwitchClientID))
	if err != nil {
		if resp.StatusCode == 404 {
			return TwitchFollowAgeErrorEmbed
		} else {
			log.Printf("GetFollowAgeEmbed: %v\n", err)
			return TwitchErrorEmbed
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GetFollowAgeEmbed: %v\n", err)
		return TwitchErrorEmbed
	}
	
	var data TwitchFollowData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("GetFollowAgeEmbed: %v\n", err)
		return TwitchErrorEmbed
	}

	if data.Error != "" {
		return TwitchFollowAgeErrorEmbed
	}

	created, _ := time.Parse("2006-01-02T15:04:05Z", data.CreatedAt)
	datestr := created.Format("Jan _2 2006")

	notiffield := &discordgo.MessageEmbedField{}
	notiffield.Name = "Notifications"
	if data.Notifications == true {
		notiffield.Value = "Notifications are enabled."
	} else {
		notiffield.Value = "Notifications are not enabled."
	}

	embed.Color = 0x6441A5
	embed.Title = "Follow Date"
	embed.Description = fmt.Sprintf("%s has followed %s since %s.", user, data.Channel.DisplayName, datestr)
	embed.Fields = append(embed.Fields, notiffield)

	return embed
}
