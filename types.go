package main

import (
	"regexp"
	"github.com/bwmarrin/discordgo"
)

/* Twitch */

type StreamStatus struct {
	Stream	*TwitchStream	`json:"stream,omitempty"`
	Links	*TwitchLinks	`json:"_links,omitempty"`

	// in case of errors:
	Error	string	`json:"error,omitempty"`
	Status	int		`json:"status,omitempty"`
	Message	string	`json:"message,omitempty"`
}

type TwitchStream struct {
	ID			uint64			`json:"_id"`
	Game		string			`json:"game"`
	Viewers		int				`json:"viewers"`
	VideoHeight	int				`json:"video_height"`
	AvgFPS		float64			`json:"average_fps"`
	Delay		int				`json:"delay"`
	CreatedAt	string			`json:"created_at"`
	IsPlaylist	bool			`json:"is_playlist"`
	StreamType	string			`json:"stream_type"`
	Previews	*TwitchPreviews	`json:"preview"`
	Channel		*TwitchChannel	`json:"channel"`
	Links		*TwitchLinks	`json:"_links"`
}

type TwitchChannel struct {
	ID			uint64	`json:"_id"`
	Mature		bool	`json:"mature"`
	Partner		bool	`json:"partner"`
	Status		string	`json:"status"`
	DisplayName	string	`json:"display_name"`
	Game		string	`json:"game"`
	Language	string	`json:"language"`
	Name		string	`json:"name"`
	CreatedAt	string	`json:"created_at"`
	UpdatedAt	string	`json:"updated_at"`
	Avatar		string	`json:"logo"`
	OfflineImg	string	`json:"video_banner"`
	URL			string	`json:"url"`
	Views		uint32	`json:"views"`
	Followers	uint32	`json:"followers"`

	// in case of error when using the channels api
	Error		string	`json:"error,omitempty"`
	Message		string	`json:"message,omitempty"`
}

type TwitchPreviews struct {
	Small		string	`json:"small"`
	Medium		string	`json:"medium"`
	Large		string	`json:"large"`
	Template	string	`json:"template"`
}

type TwitchLinks struct {
	Self	string	`json:"self"`
	Channel	string	`json:"channel,omitempty"`
}

/* Commands */

type MsgHandler func(*discordgo.MessageCreate, *discordgo.Session)

type Command struct {
	Pattern *regexp.Regexp
	Name string
	Args []CommandArg
	Examples []string
	Description string
	Aliases []string
	Handler MsgHandler
	Category string // if "" the command won't be listed in help menu
	AdminOnly bool
}

type CommandArg struct {
	Title string
	Required bool
}
