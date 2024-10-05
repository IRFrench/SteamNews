package discord

import (
	"crypto/tls"
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
	// Discord's CA certificates aren't in alpine??
	return DiscordClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		auth: fmt.Sprintf("Bot %s", botToken),
	}
}
