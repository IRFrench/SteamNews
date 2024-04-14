package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DmBody struct {
	RecipientId int `json:"recipient_id,omitempty"`
}

type Channel struct {
	Id   string `json:"id,omitempty"`
	Type int    `json:"type,omitempty"`
}

func (d *DiscordClient) CreateDmChannel() (Channel, error) {
	discordUrl := fmt.Sprintf("%s/v%d/users/@me/channels", discordUrl, discordApiVersion)

	parsedBody, err := json.Marshal(DmBody{RecipientId: d.userId})
	if err != nil {
		return Channel{}, err
	}

	request, err := http.NewRequest(http.MethodPost, discordUrl, bytes.NewBuffer(parsedBody))
	if err != nil {
		return Channel{}, err
	}
	defer request.Body.Close()

	response, err := d.sendRequest(request)
	if err != nil {
		return Channel{}, err
	}

	var jsonResponse Channel
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return Channel{}, err
	}

	return jsonResponse, err
}
