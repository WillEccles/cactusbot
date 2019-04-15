package main

import (
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
)

type Configuration struct {
	DiscordToken string
	DiscordClientID string
	DebugChannel string
	AdminIDs []string
	ControllerID string
}

func LoadConfig() Configuration {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		WriteNewConfig()
		return Configuration{}
	}
	
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println(err)
	}
	return conf
}

func WriteNewConfig() {
	file, _ := json.MarshalIndent(Configuration{}, "", "\t")
	_ = ioutil.WriteFile("config.json", file, 0644)
}

func (c *Configuration) IsAdmin(id string) bool {
	for _, i := range(c.AdminIDs) {
		if i == id {
			return true
		}
	}
	return false
}
