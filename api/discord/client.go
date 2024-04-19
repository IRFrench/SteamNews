package discord

import (
	"fmt"
	"net/http"
)

const (
	DISCORD_API_VERSION = 10
	DISCORD_URL         = "https://discord.com/api"
)

type DiscordClient struct {
	client *http.Client
	auth   string
}

func (d *DiscordClient) sendRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add(
		"Authorization",
		d.auth,
	)
	request.Header.Add(
		"User-Agent",
		fmt.Sprintf("DiscordBot (%s, %d)", DISCORD_URL, DISCORD_API_VERSION),
	)
	request.Header.Add(
		"Content-Type",
		"application/json",
	)
	return d.client.Do(request)
}

func NewDiscordClient(botToken string) DiscordClient {
	return DiscordClient{
		client: &http.Client{},
		auth:   fmt.Sprintf("Bot %s", botToken),
	}
}
