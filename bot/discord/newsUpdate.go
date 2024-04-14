package discord

import "time"

type Game struct {
	Name string
	News []ShortArticle
}

type ShortArticle struct {
	Title    string
	Url      string
	Author   string
	Contents string
	Date     time.Time
}

type Message struct {
	Content string `json:"content,omitempty"`
}

// func (d *DiscordClient) SendNewsMessage(news []Game) error {

// }
