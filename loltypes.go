package main

import "sync"

// use this to access all the data about league
type LeagueHelper struct {
    Token   string
    Version string

    ChampionData *ChampionFile // current champion data file's contents

    Lock sync.Mutex // used to lock this struct while updating it
}

/* Data Types */

type LeagueStatus struct {
    StatusCode  int     `json:"status_code,omitempty"`
    Message     string  `json:"message,omitempty"`
}

type Summoner struct {
    ProfileIconID   int     `json:"profileIconId,omitempty"`
    Name            string  `json:"name,omitempty"`
    PUUID           string  `json:"puuid,omitempty"`
    Level           int64   `json:"summonerLevel,omitempty"`
    RevisionDate    int64   `json:"revisionDate,omitempty"`
    ID              string  `json:"id,omitempty"`
    AccountID       string  `json:"accountId,omitempty"`
    
    Status  *LeagueStatus   `json:"status,omitempty"`
}

type ChampionMasteries []ChampionMasteryDTO

type ChampionMasteryDTO struct {
    ChestGranted            bool    `json:"chestGranted,omitempty"`
    ChampionLevel           int     `json:"championLevel,omitempty"`
    ChampionPoints          int     `json:"championPoints,omitempty"`
    ChampionID              int     `json:"championId,omitempty"`
    PointsUntilNextLevel    int64   `json:"championPointsUntilNextLevel,omitempty"`
    LastPlayTime            int64   `json:"lastPlayTime,omitempty"`
    TokensEarned            int     `json:"tokensEarned,omitempty"`
    PointsSinceLastLevel    int64   `json:"championPointsSinceLastLevel,omitempty"`
    SummonerID              string  `json:"summonerId,omitempty"`

    Status  *LeagueStatus   `json:"status,omitempty"`
}

type GenericLeagueError struct {
    Status *LeagueStatus    `json:"status,omitempty"`
}

type ChampionFile struct {
    Format  string  `json:"format,omitempty"`
    Version string  `json:"version,omitempty"`

    Data    map[string]*ChampionDTO `json:"data,omitempty"`
}

// compatible with champion.json, championFull.json, and also <championname>.json
type ChampionDTO struct {
    Version     string      `json:"version,omitempty"`
    ID          string      `json:"id,omitempty"`
    Key         string      `json:"key,omitempty"`
    Name        string      `json:"name,omitempty"`
    Title       string      `json:"title,omitempty"`
    Blurb       string      `json:"blurb,omitempty"`
    Lore        string      `json:"lore,omitempty"`
    Tags        []string    `json:"tags,omitempty"`
    Resource    string      `json:"partype,omitempty"`


    Info    *ChampionInfoDTO    `json:"info,omitempty"`
    Image   *ImageDTO           `json:"image,omitempty"`
    Stats   *ChampionStatsDTO   `json:"stats,omitempty"`
    Spells  []*ChampionSpellDTO `json:"spells,omitempty"`
    Passive *ChampionPassiveDTO `json:"passive,omitempty"`
}

type ChampionInfoDTO struct {
    Attack      int `json:"attack,omitempty"`
    Defense     int `json:"defense,omitempty"`
    Magic       int `json:"magic,omitempty"`
    Difficulty  int `json:"difficulty,omitempty"`
}

type ImageDTO struct {
    Full    string  `json:"full,omitempty"`
    Sprite  string  `json:"sprite,omitempty"`
    Group   string  `json:"group,omitempty"`
    PosX    int     `json:"x,omitempty"`
    PosY    int     `json:"y,omitempty"`
    Width   int     `json:"w,omitempty"`
    Height  int     `json:"h,omitempty"`
}

type ChampionStatsDTO struct {
    HP                      float32 `json:"hp,omitempty"`
    HPPerLevel              float32 `json:"hpperlevel,omitempty"`
    MP                      float32 `json:"mp,omitempty"`
    MPPerLevel              float32 `json:"mpperlevel,omitempty"`
    MoveSpeed               float32 `json:"movespeed,omitempty"`
    Armor                   float32 `json:"armor,omitempty"`
    ArmorPerLevel           float32 `json:"armorperlevel,omitempty"`
    MagicResist             float32 `json:"spellblock,omitempty"`
    MagicResistPerLevel     float32 `json:"spellblockperlevel,omitempty"`
    AttackRange             float32 `json:"attackrange,omitempty"`
    HPRegen                 float32 `json:"hpregen,omitempty"`
    HPRegenPerLevel         float32 `json:"hpregenperlevel,omitempty"`
    MPRegen                 float32 `json:"mpregen,omitempty"`
    MPRegenPerLevel         float32 `json:"mpregenperlevel,omitempty"`
    Crit                    float32 `json:"crit,omitempty"`
    CritPerLevel            float32 `json:"critperlevel,omitempty"`
    AttackDamage            float32 `json:"attackdamage,omitempty"`
    AttackDamagePerLevel    float32 `json:"attackdamageperlevel,omitempty"`
    AttackSpeed             float32 `json:"attackspeed,omitempty"`
    AttackSpeedPerLevel     float32 `json:"attackspeedperlevel,omitempty"`
}

type ChampionSpellDTO struct {
    ID          string  `json:"id,omitempty"`
    Name        string  `json:"name,omitempty"`
    Description string  `json:"description,omitempty"`
    // will be "No Cost" if no cost
    Cost        string  `json:"resource,omitempty"`
    // will be "No Cost" if no cost
    CostType    string  `json:"costType,omitempty"`
    // will be "-1" if no ammo
    MaxAmmo     string  `json:"maxammo,omitempty"`

    Image   *ImageDTO   `json:"image,omitempty"`
}

type ChampionPassiveDTO struct {
    Name        string  `json:"name,omitempty"`
    Description string  `json:"description,omitempty"`

    Image   *ImageDTO   `json:"image,omitempty"`
}

type ServerStatus struct {
    Name        string      `json:"name,omitempty"`
    RegionTag   string      `json:"region_tag,omitempty"`
    Hostname    string      `json:"hostname,omitempty"`
    Slug        string      `json:"slug,omitempty"`
    Locales     []string    `json:"locales,omitempty"`

    Services    []*Service  `json:"services,omitempty"`
    Status  *LeagueStatus   `json:"status,omitempty"`
}

type Service struct {
    Status  string  `json:"status,omitempty"`
    Name    string  `json:"name,omitempty"`
    Slug    string  `json:"slug,omitempty"`

    Incidents   []*Incident `json:"incidents,omitempty"`
}

type Incident struct {
    Active      bool    `json:"active,omitempty"`
    CreatedAt   string  `json:"created_at,omitempty"`
    ID          uint64  `json:"id,omitempty"`

    Updates []*Message  `json:"updates,omitempty"`
}

type Message struct {
    Severity    string  `json:"severity,omitempty"`
    Author      string  `json:"author,omitempty"`
    CreatedAt   string  `json:"created_at,omitempty"`
    UpdatedAt   string  `json:"updated_at,omitempty"`
    Content     string  `json:"content,omitempty"`
    ID          string  `json:"id,omitempty"`

    Translations    []*Translation  `json:"translations,omitempty"`
}

type Translation struct {
    Locale      string  `json:"locale,omitempty"`
    Content     string  `json:"content,omitempty"`
    UpdatedAt   string  `json:"updated_at,omitempty"`
}
