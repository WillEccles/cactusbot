package main

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
)

const StreamStatusURL = "https://api.twitch.tv/kraken/streams/%s?client_id=%s"

func GetStreamStatusEmbed(stream string) *discordgo.MessageEmbed {
	//embed := &discordgo.MessageEmbed{}
	
	resp, err := http.Get(fmt.Sprintf(StreamStatusURL, stream, Config.TwitchClientID))
	if err != nil {
		// TODO do error things here
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// TODO do error things here
		return nil
	}
	
	var status StreamStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		// TODO do error things here
		return nil
	}
	
	return nil
}
