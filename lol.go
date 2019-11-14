package main

import (
	"net/http"
    "net/url"
	"encoding/json"
	"fmt"
	"io/ioutil"
    "log"
    "strings"
    "strconv"
    "time"
    "os"

	"github.com/bwmarrin/discordgo"
)

const (
    // TODO make this non-NA-specific
    API_URL = "https://na1.api.riotgames.com"
    // takes summoner name, URL encoded (obviously)
    SUMMONER = "/lol/summoner/v4/summoners/by-name/%s"
    // takes summonerID
    ALL_CHAMPION_MASTERY = "/lol/champion-mastery/v4/champion-masteries/by-summoner/%s"
    // takes championID; tack onto end of ALL_CHAMPION_MASTERY
    BY_CHAMPION = "/by-champion/%s"
    // takes summonerID
    MASTERY_SCORE = "/lol/champion-mastery/v4/scores/by-summoner/%v"
    // takes a version and profile icon ID
    PROFILE_ICON = "http://ddragon.leagueoflegends.com/cdn/%v/img/profileicon/%v.png"
    // takes version, group, image name
    IMAGE_URL_PATTERN = "http://ddragon.leagueoflegends.com/cdn/%v/img/%v/%v"
    // gets the versions
    VERSIONS_URL = "https://ddragon.leagueoflegends.com/api/versions.json"
    // gets championFull.json
    CHAMPION_FULL = "http://ddragon.leagueoflegends.com/cdn/%v/data/en_US/championFull.json"
)

// returns true if all is well; false if not
func (helper *LeagueHelper) Init(token string) bool {
    helper.Token = token

    downloaded := false
    if _, err := os.Stat("championFull.json"); os.IsNotExist(err) {
        log.Println("Champion data not found; getting latest...")
        ver := getLatestVersion()
        if ver != "" {
            log.Printf("Current version: %s\n", ver)
            helper.Version = ver
            success := downloadChampionData(ver)
            if !success {
                return false
            }
            downloaded = true
            log.Println("Latest champion data acquired.")
        } else {
            return false
        }
    }

    helper.ChampionData = loadChampionsFile()
    if helper.ChampionData != nil {
        helper.Version = helper.ChampionData.Version
    } else {
        return false
    }
    
    if !downloaded {
        helper.UpdateData()
    }

    return true
}

// assumes the file exists; check first! returns nil if error
func loadChampionsFile() *ChampionFile {
    cfile := &ChampionFile{}

    fcontents, err := ioutil.ReadFile("championFull.json")
	if err != nil {
        log.Printf("Error loading championFull.json:\n%v\n", err)
		return nil
	}
	
    err = json.Unmarshal(fcontents, cfile)
    if err != nil {
        log.Printf("Error parsing championFull.json:\n%v\nPlease delete the file and run the bot again.\n", err)
        return nil
    }
	
	return cfile
}

func getLatestVersion() string {
    var vers []string

    resp, err := http.Get(VERSIONS_URL)
    if err != nil {
        if resp.StatusCode != 200 && resp.StatusCode != 404 {
            log.Printf("Error in getLatestVersion:\n%v\n", err)
            return ""
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in getLatestVersion:\n%v\n", err)
        return ""
    }

    err = json.Unmarshal(body, &vers)
    if err != nil {
        log.Printf("Error in getLatestVersion:\n%v\n", err)
        return ""
    }

    return vers[0]
}

func downloadChampionData(version string) bool {
    resp, err := http.Get(fmt.Sprintf(CHAMPION_FULL, version))
    if err != nil {
        log.Printf("Error in downloadChampionData:\n%v\n", err)
        return false
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in downloadChampionData:\n%v\n", err)
        return false
    }

    err = ioutil.WriteFile("championFull.json", body, 0644)
    if err != nil {
        log.Printf("Error in downloadChampionData:\n%v\n", err)
        return false
    }

    return true
}

// returns true if updated, string is error or ""
func (helper *LeagueHelper) UpdateData() (bool, string) {
    helper.Lock.Lock()
    defer helper.Lock.Unlock()
    latestver := getLatestVersion()
    if latestver == "" {
        return false, "Error getting latest version"
    }

    if helper.Version != latestver {
        log.Printf("Updating League data from %v to %v...\n", helper.Version, latestver)
        // do actual update
        // get the latest champion data
        success := downloadChampionData(latestver)
        if !success {
            log.Println("Failed to update League data.")
            return false, "Failed to update League data"
        }
        // load the latest data
        cfile := loadChampionsFile()
        if cfile == nil {
            log.Println("Failed to update League data")
            return false, "Failed to update League data"
        }
        // set the latest data into the helper
        helper.ChampionData = cfile
        log.Printf("Successfully updated League data to %v\n", latestver)
        return true, ""
    } else {
        log.Printf("League data is up-to-date. (%v)\n", helper.Version)
        return false, ""
    }

    return true, ""
}

// should only be run in a separate goroutine
func (helper *LeagueHelper) UpdateRoutine() {
    for {
        time.Sleep(
    }
}

func (i *ImageDTO) GetURL(version string) string {
    return fmt.Sprintf(IMAGE_URL_PATTERN, version, i.Group, i.Full)
}

func (s *Summoner) GetIconURL(version string) string {
    return fmt.Sprintf(PROFILE_ICON, version, s.ProfileIconID)
}

func MakeErrorEmbed(err string) *discordgo.MessageEmbed {
    return &discordgo.MessageEmbed{
        Color: 0xCC0000,
        Description: err,
    }
}

func (helper LeagueHelper) GetSummoner(summonername string) (*Summoner, string) {
    helper.Lock.Lock()
    defer helper.Lock.Unlock()
    requrl := fmt.Sprintf(API_URL+SUMMONER, url.QueryEscape(summonername)) + "?api_key=" + helper.Token
    requrl = strings.ReplaceAll(requrl, "+", "%20")

    resp, err := http.Get(requrl)
    if err != nil {
        if resp.StatusCode != 200 && resp.StatusCode != 404 {
            log.Printf("Error in GetSummoner:\n%v\n", err)
            return nil, "Error retrieving data from API"
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in GetSummoner:\n%v\n", err)
        return nil, "Error reading API data"
    }
    
    summ := &Summoner{}
    err = json.Unmarshal(body, summ)
    if err != nil {
        log.Printf("Error in GetSummoner:\n%v\n", err)
        return nil, "Error parsing API data"
    }

    return summ, ""
}

func (helper LeagueHelper) GetMasteryScore(summonerID string) (int, string) {
    helper.Lock.Lock()
    defer helper.Lock.Unlock()
    requrl := fmt.Sprintf(API_URL+MASTERY_SCORE + "?api_key=" + helper.Token, summonerID)
    
    waserr := false
    resp, err := http.Get(requrl)
    if err != nil {
        if resp.StatusCode != 200 && resp.StatusCode != 404 {
            log.Printf("Error in GetMasteryScore:\n%v\n", err)
            return -1, "Error retrieving data from API"
        } else {
            waserr = true
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in GetMasteryScore:\n%v\n", err)
        return -1, "Error reading API data"
    }

    if !waserr {
        mastery, err := strconv.Atoi(string(body))
        if err != nil {
            log.Printf("Error in GetMasteryScore:\n%v\n", err)
            return -1, "Error parsing API data"
        }

        return mastery, ""
    } else {
        var e GenericLeagueError
        err = json.Unmarshal(body, &e)
        if err != nil {
            log.Printf("Error in GetMasteryScore:\n%v\n", err)
            return -1, "Error parsing error data from API"
        }
        log.Printf("Error in GetMasteryScore:\n%v\n", e.Status.Message)
        return -1, e.Status.Message
    }
}

func (helper LeagueHelper) GetSummonerEmbed(summonername string) *discordgo.MessageEmbed {
    embed := &discordgo.MessageEmbed{}

    summoner, err := helper.GetSummoner(summonername)
    if err != "" {
        log.Printf("Error in GetSummonerEmbed:\n%v\n", err)
        return MakeErrorEmbed(err)
    }
    if summoner.Status != nil {
        log.Printf("Error in GetSummonerEmbed:\n%v\n", summoner.Status.Message)
        return MakeErrorEmbed(summoner.Status.Message)
    }

    mastery, err := helper.GetMasteryScore(summoner.ID)
    if err != "" {
        return MakeErrorEmbed(err)
    }

    updatetime := time.Unix(summoner.RevisionDate / 1000, 0)
    updatestamp := updatetime.Format(time.RFC1123)

    embed.Color = 0xD13739
    embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
        URL: summoner.GetIconURL(helper.Version),
    }
    embed.Title = "Summoner: " + summoner.Name
    embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
        Name: "Level",
        Value: fmt.Sprintf("%v", summoner.Level),
    }, &discordgo.MessageEmbedField{
        Name: "Mastery",
        Value: strconv.Itoa(mastery),
    }, &discordgo.MessageEmbedField{
        Name: "Updated",
        Value: updatestamp,
    })

    return embed
}
