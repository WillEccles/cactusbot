package main

import (
    "encoding/json"
    "log"
    "io/ioutil"
)

// omited from types.go, makes more sense to be in here
type Configuration struct {
    DiscordToken    string      `json:",omitempty"`
    DiscordClientID string      `json:",omitempty"`
    DebugChannel    string      `json:",omitempty"`
    AdminIDs        []string    `json:",omitempty"`
    ControllerID    string      `json:",omitempty"`
    LogChannel      string      `json:",omitempty"`
    LogWebhookID    string      `json:",omitempty"`
    LogWebhookToken string      `json:",omitempty"`
    LeagueToken     string      `json:",omitempty"`
}

func LoadConfig() Configuration {
    var conf Configuration
    fcontents, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Printf("Error loading config:\n%v\n", err)
        return Configuration{}
    }
    
    err = json.Unmarshal(fcontents, &conf)
    if err != nil {
        log.Printf("Error parsing json config:\n%v\n", err)
        return Configuration{}
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
