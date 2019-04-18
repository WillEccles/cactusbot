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

const StreamStatusURL = "https://api.twitch.tv/kraken/streams/%s?client_id=%s"
const ChannelInfoURL = "https://api.twitch.tv/kraken/channels/%s?client_id=%s"
var TwitchErrorEmbed = &discordgo.MessageEmbed{
	Title: "Error",
	Description: "Error getting info for the stream! Try again later.",
	Color: 0xCE0000,
}
var TwitchNoSuchChannelEmbed = &discordgo.MessageEmbed{
	Title: "Error",
	Description: "That channel does not exist.",
	Color: 0xCE0000,
}

func GetStreamStatusEmbed(stream string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}
	
	resp, err := http.Get(fmt.Sprintf(StreamStatusURL, stream, Config.TwitchClientID))
	if err != nil {
		log.Printf("GetStreamStatusEmbed: %v\n", err)
		return TwitchErrorEmbed
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
	
	resp, err := http.Get(fmt.Sprintf(StreamStatusURL, stream, Config.TwitchClientID))
	if err != nil {
		log.Printf("GetChannelInfoEmbed: %v\n", err)
		return TwitchErrorEmbed
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
