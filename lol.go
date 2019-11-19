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
    "regexp"

    "github.com/bwmarrin/discordgo"
)

const (
    // TODO make this non-NA-specific
    API_URL = "https://na1.api.riotgames.com"
    // takes summoner name, URL encoded (obviously)
    SUMMONER = "/lol/summoner/v4/summoners/by-name/%s"
    // takes summonerID
    ALL_CHAMPION_MASTERY = "/lol/champion-mastery/v4/champion-masteries/by-summoner/%v"
    // takes championID; tack onto end of ALL_CHAMPION_MASTERY
    BY_CHAMPION = "/by-champion/%v"
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
    // gets server status
    SERVER_STATUS = "/lol/status/v3/shard-data"
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
        _, err := helper.UpdateData()
        if err != "" {
            log.Println("Error in helper.Init:\n%v\n", err)
            return false
        }
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

    nmap := make(map[string]*ChampionDTO)
    for champ, champdata := range(cfile.Data) {
        nmap[strings.ToLower(strings.ReplaceAll(champ, "'", ""))] = champdata
    }
    cfile.Data = nmap
    
    return cfile
}

func getLatestVersion() string {
    var vers []string

    resp, err := http.Get(VERSIONS_URL)
    if err != nil {
        if resp != nil {
            if resp.StatusCode != 200 && resp.StatusCode != 404 {
                log.Printf("Error in getLatestVersion:\n%v\n", err)
                return ""
            }
        } else {
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
        log.Printf("League data is up-to-date (version %v)\n", helper.Version)
        return false, ""
    }

    return true, ""
}

// should only be run in a separate goroutine
func (helper *LeagueHelper) UpdateRoutine() {
    for {
        time.Sleep(12 * time.Hour)
        helper.UpdateData()
    }
}

func (i *ImageDTO) GetURL(version string) string {
    return fmt.Sprintf(IMAGE_URL_PATTERN, version, i.Group, i.Full)
}

func (s *Summoner) GetIconURL(version string) string {
    return fmt.Sprintf(PROFILE_ICON, version, s.ProfileIconID)
}

func sanitizeChampionName(champname string) string {
    return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(champname), " ", ""), "'", "")
}

// returns -1 if not found
func (helper *LeagueHelper) getChampionIDByName(name string) int {
    sname := sanitizeChampionName(name)
    for k, v := range(helper.ChampionData.Data) {
        if k == sname {
            idval, _ := strconv.Atoi(v.Key)
            return idval
        }
    }
    return -1
}

// returns "" if not found
func (helper *LeagueHelper) getChampionNameByID(id int) string {
    for _, v := range(helper.ChampionData.Data) {
        key, _ := strconv.Atoi(v.Key)
        if key == id {
            return v.Name
        }
    }
    return ""
}

func MakeErrorEmbed(err string) *discordgo.MessageEmbed {
    return &discordgo.MessageEmbed{
        Color: 0xCC0000,
        Description: err,
    }
}

func (helper *LeagueHelper) getStatus() (*ServerStatus, string) {
    requrl := API_URL + SERVER_STATUS + "?api_key=" + helper.Token

    resp, err := http.Get(requrl)
    if err != nil {
        if resp.StatusCode != 200 && resp.StatusCode != 404 {
            log.Printf("Error in getStatus:\n%v\n", err)
            return nil, "Error retrieving data from API"
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in getStatus:\n%v\n", err)
        return nil, "Error reading data from API"
    }

    status := &ServerStatus{}
    err = json.Unmarshal(body, &status)
    if err != nil {
        log.Printf("Error in getStatus:\n%v\n", err)
        return nil, "Error parsing data from the API"
    }

    return status, ""
}

func (helper *LeagueHelper) GetStatusEmbed() *discordgo.MessageEmbed {
    embed := &discordgo.MessageEmbed{}

    status, err := helper.getStatus()
    if err != "" {
        return MakeErrorEmbed(err)
    }

    if status.Status != nil {
        return MakeErrorEmbed(status.Status.Message)
    }

    embed.Title = "Server Status (NA1)"
    embed.Color = 0xD13739
    for _, service := range(status.Services) {
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: service.Name,
            Value: strings.Title(service.Status),
        })
    }

    return embed
}

func (helper LeagueHelper) GetSummoner(summonername string) (*Summoner, string) {
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
    requrl := fmt.Sprintf(API_URL+MASTERY_SCORE + "?api_key=" + helper.Token, summonerID)
    
    waserr := false
    resp, err := http.Get(requrl)
    if err != nil {
        if resp != nil {
            if resp.StatusCode != 200 && resp.StatusCode != 404 {
                log.Printf("Error in GetMasteryScore:\n%v\n", err)
                return -1, "Error retrieving data from API"
            } else {
                waserr = true
            }
        } else {
            log.Printf("Error in GetMasteryScore:\n%v\n", err)
            return -1, "Error retrieving data from API"
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

func (helper LeagueHelper) GetSummonerMasteries(summonerID string) (ChampionMasteries, string) {
    requrl := fmt.Sprintf(API_URL+ALL_CHAMPION_MASTERY+"?api_key="+helper.Token, summonerID)
    
    waserr := false
    resp, err := http.Get(requrl)
    if err != nil {
        if resp != nil {
            if resp.StatusCode != 200 && resp.StatusCode != 404 {
                log.Printf("Error in GetSummonerMasteries:\n%v\n", err)
                return nil, "Error retrieving data from API"
            } else {
                waserr = true
            }
        } else {
            log.Printf("Error in GetSummonerMasteries:\n%v\n", err)
            return nil, "Error retrieving data from API"
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in GetSummonerMasteries:\n%v\n", err)
        return nil, "Error reading API data"
    }

    if waserr {
        lerr := &GenericLeagueError{}
        err = json.Unmarshal(body, lerr)
        if err != nil {
            log.Println("Error in GetSummonerMasteries:\n%v\n", err)
            return nil, "Error parsing API data"
        }
        return nil, lerr.Status.Message
    }

    var m ChampionMasteries

    err = json.Unmarshal(body, &m)
    if err != nil {
        log.Println("Error in GetSummonerMasteries:\n%v\n", err)
        return nil, "Error parsing API data"
    }

    return m, ""
    
}

func (helper LeagueHelper) GetSummonerMasteryForChampion(summonerID, champ string) (*ChampionMasteryDTO, string) {
    cid := helper.getChampionIDByName(champ)
    if cid == -1 {
        return nil, "Champion not found: " + champ
    }

    requrl := fmt.Sprintf(API_URL+ALL_CHAMPION_MASTERY+BY_CHAMPION + "?api_key=" + helper.Token, summonerID, cid)
   
    waserr := false
    resp, err := http.Get(requrl)
    if err != nil {
        if resp != nil {
            if resp.StatusCode != 200 && resp.StatusCode != 404 {
                log.Printf("Error in GetSummonerMasteryForChampion:\n%v\n", err)
                return nil, "Error retrieving data from API"
            } else {
                waserr = true
            }
        } else {
            log.Printf("Error in GetSummonerMasteryForChampion:\n%v\n", err)
            return nil, "Error retrieving data from API"
        }
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error in GetSummonerMasteryForChampion\n%v\n", err)
        return nil, "Error reading API data"
    }

    if waserr {
        lerr := &GenericLeagueError{}
        err = json.Unmarshal(body, lerr)
        if err != nil {
            log.Println("Error in GetSummonerMasteryForChampion:\n%v\n", err)
            return nil, "Error parsing API data"
        }
        return nil, lerr.Status.Message
    }

    mastery := &ChampionMasteryDTO{}

    err = json.Unmarshal(body, mastery)
    if err != nil {
        log.Println("Error in GetSummonerMasteryForChampion:\n%v\n", err)
        return nil, "Error parsing API data"
    }

    return mastery, ""
}

func (helper *LeagueHelper) GetSummonerEmbed(summonername string) *discordgo.MessageEmbed {
    helper.Lock.Lock()
    defer helper.Lock.Unlock()
    embed := &discordgo.MessageEmbed{}

    summoner, err := helper.GetSummoner(summonername)
    if err != "" {
        log.Printf("Error in GetSummonerEmbed:\n%v\n", err)
        return MakeErrorEmbed(err)
    }
    if summoner.Status != nil {
        //log.Printf("Error in GetSummonerEmbed:\n%v\n", summoner.Status.Message)
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

func (helper *LeagueHelper) GetSummonerMasteriesEmbed(summonername string) *discordgo.MessageEmbed {
    helper.Lock.Lock()
    defer helper.Lock.Unlock()
    embed := &discordgo.MessageEmbed{}

    summoner, err := helper.GetSummoner(summonername)
    if err != "" {
        log.Printf("Error in GetSummonerMasteriesEmbed:\n%v\n", err)
        return MakeErrorEmbed(err)
    }
    if summoner.Status != nil {
        //log.Printf("Error in GetSummonerEmbed:\n%v\n", summoner.Status.Message)
        return MakeErrorEmbed(summoner.Status.Message)
    }
    
    mastery, err := helper.GetMasteryScore(summoner.ID)
    if err != "" {
        return MakeErrorEmbed(err)
    }

    masteries, err := helper.GetSummonerMasteries(summoner.ID)
    if err != "" {
        return MakeErrorEmbed(err)
    }

    totalpoints := 0
    for _, m := range masteries {
        totalpoints += m.ChampionPoints
    }

    embed.Color = 0xD13739
    embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
        URL: summoner.GetIconURL(helper.Version),
    }
    embed.Title = "Summoner Masteries: " + summoner.Name
    embed.Description = fmt.Sprintf("**Mastery level: %v**\nTotal mastery points: %v", mastery, totalpoints)
    for i := 0; i < 5; i++ {
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: helper.getChampionNameByID(masteries[i].ChampionID),
            Value: fmt.Sprintf("Level %v, %vpts", masteries[i].ChampionLevel, masteries[i].ChampionPoints),
        })
    }

    return embed
}

func (helper *LeagueHelper) GetSummonerMasteryEmbed(summonername, champname string) *discordgo.MessageEmbed {
    helper.Lock.Lock()
    defer helper.Lock.Unlock()
    embed := &discordgo.MessageEmbed{}

    summoner, err := helper.GetSummoner(summonername)
    if err != "" {
        log.Printf("Error in GetSummonerMasteriesEmbed:\n%v\n", err)
        return MakeErrorEmbed(err)
    }
    if summoner.Status != nil {
        //log.Printf("Error in GetSummonerEmbed:\n%v\n", summoner.Status.Message)
        return MakeErrorEmbed(summoner.Status.Message)
    }

    mastery, err := helper.GetSummonerMasteryForChampion(summoner.ID, champname)
    if err != "" {
        return MakeErrorEmbed(err)
    }

    // get details for champ
    champ := helper.ChampionData.Data[sanitizeChampionName(champname)]

    lastplaytime := time.Unix(mastery.LastPlayTime / 1000, 0)
    lastplaytimestamp := lastplaytime.Format(time.RFC1123)

    embed.Color = 0xD13739
    embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
        URL: champ.Image.GetURL(helper.Version),
    }
    embed.Title = "Champion Mastery: " + champ.Name
    embed.Description = "For summoner " + summoner.Name
    embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
        Name: "Level",
        Value: fmt.Sprintf("%v", mastery.ChampionLevel),
    }, &discordgo.MessageEmbedField{
        Name: "Points",
        Value: fmt.Sprintf("%v", mastery.ChampionPoints),
    }, &discordgo.MessageEmbedField{
        Name: "Tokens",
        Value: fmt.Sprintf("%v", mastery.TokensEarned),
    }, &discordgo.MessageEmbedField{
        Name: "Last Played",
        Value: lastplaytimestamp,
    })

    return embed
    
}

var brtag = regexp.MustCompile(`(?i)(<br/?>)+`)
var htmltag = regexp.MustCompile(`(?i)<[^<]+>`)
func sanitizeDescription(desc string) string {
    return htmltag.ReplaceAllString(brtag.ReplaceAllString(desc, "\n"), "*")
}

func (helper *LeagueHelper) GetChampionEmbed(champname string) *discordgo.MessageEmbed {
    embed := &discordgo.MessageEmbed{}
    
    cdata, found := helper.ChampionData.Data[sanitizeChampionName(champname)]
    if !found {
        return MakeErrorEmbed("Error: Champion not found")
    }

    embed.Color = 0xD13739
    embed.Title = "Champion: " + cdata.Name
    embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
        URL: cdata.Image.GetURL(helper.Version),
    }
    embed.Description = fmt.Sprintf("**%v**\n*", strings.Title(cdata.Title))
    for i, t := range(cdata.Tags) {
        if i != 0 {
            embed.Description += ", "
        }
        embed.Description += t
    }
    embed.Description += "*"
    embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
        Name: "Lore",
        Value: cdata.Lore,
    }, &discordgo.MessageEmbedField{
        Name: "Resource",
        Value: cdata.Resource,
    })

    spellLabels := "QWER"
    for i, spell := range(cdata.Spells) {
        f := &discordgo.MessageEmbedField{}
        f.Name = fmt.Sprintf("%c: %v", spellLabels[i], spell.Name)
        f.Value = sanitizeDescription(spell.Description)
        embed.Fields = append(embed.Fields, f)
    }

    embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
        Name: "Passive: " + cdata.Passive.Name,
        Value: sanitizeDescription(cdata.Passive.Description),
    })

    return embed
}
