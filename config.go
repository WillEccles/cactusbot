package main

import (
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
)

// omited from types.go, makes more sense to be in here
type Configuration struct {
    DiscordToken	string
    DiscordClientID	string
    DebugChannel	string
    AdminIDs		[]string
    ControllerID	string
    LogChannel      string
    LogWebhookID    string
    LogWebhookToken string
}

func LoadConfig() Configuration {
	file, err := os.Open("config.json")
	if err != nil {
        log.Printf("Error loading config:\n%v\n", err)
		return Configuration{}
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println(err)
	}

	return conf
}

func WriteConfig(conf Configuration) {
	// write the config to file, which ensures all config files stay up to date when new options are added
	file, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Println("WriteConfig: Error marshalling json for writing.")
		return
	}
	err = ioutil.WriteFile("config.json", file, 0644)
	if err != nil {
		log.Printf("WriteConfig: Error writing config to file! Config:\n%v\n", string(file))
	}
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
