package discord

import (
	"fmt"
	"net/http"
)

const (
	discordApiVersion = 10
	discordUrl        = "https://discord.com/api"
)

type DiscordClient struct {
	client *http.Client
	auth   string
	userId int
}

func (d *DiscordClient) sendRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add(
		"Authorization",
		d.auth,
	)
	request.Header.Add(
		"User-Agent",
		fmt.Sprintf("DiscordBot (%s, %d)", discordUrl, discordApiVersion),
	)
	request.Header.Add(
		"Content-Type",
		"application/json",
	)
	return d.client.Do(request)
}

func NewDiscordClient(botToken string, userId int) DiscordClient {
	return DiscordClient{
		client: &http.Client{},
		auth:   fmt.Sprintf("Bot %s", botToken),
		userId: userId,
	}
}
