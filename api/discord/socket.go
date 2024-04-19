package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SocketResponse struct {
	Url               string       `json:"url,omitempty"`
	Shards            int          `json:"shards,omitempty"`
	SessionStartLimit SessionStart `json:"session_start_limit,omitempty"`
	Message           string       `json:"message,omitempty"`
	Code              int          `json:"code,omitempty"`
}

type SessionStart struct {
	Total          int `json:"total,omitempty"`
	Remaining      int `json:"remaining,omitempty"`
	ResetAfter     int `json:"reset_after,omitempty"`
	MaxConcurrency int `json:"max_concurrency,omitempty"`
}

func (d *DiscordClient) CollectSocket() (SocketResponse, error) {
	discordUrl := fmt.Sprintf("%s/v%d/gateway/bot", DISCORD_URL, DISCORD_API_VERSION)

	request, err := http.NewRequest(http.MethodGet, discordUrl, nil)
	if err != nil {
		return SocketResponse{}, err
	}

	response, err := d.sendRequest(request)
	if err != nil {
		return SocketResponse{}, err
	}
	defer response.Body.Close()

	var jsonResponse SocketResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return SocketResponse{}, err
	}

	if response.StatusCode != 200 {
		return SocketResponse{}, fmt.Errorf("%s (%d)", jsonResponse.Message, jsonResponse.Code)
	}

	return jsonResponse, nil
}
