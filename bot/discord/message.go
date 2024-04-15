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

type MessageRecipient struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Id            string `json:"id"`
	Avatar        string `json:"avatar"`
}

type Message struct {
	LastMessageID string             `json:"last_message_id,omitempty"`
	Type          int                `json:"type,omitempty"`
	Id            string             `json:"id,omitempty"`
	Recipients    []MessageRecipient `json:"recipients,omitempty"`
	Message       string             `json:"message,omitempty"`
	Code          int                `json:"code,omitempty"`
}

func (d *DiscordClient) CreateDmChannel(userId int) (Message, error) {
	discordUrl := fmt.Sprintf("%s/v%d/users/@me/channels", discordUrl, discordApiVersion)

	parsedBody, err := json.Marshal(DmBody{RecipientId: userId})
	if err != nil {
		return Message{}, err
	}

	request, err := http.NewRequest(http.MethodPost, discordUrl, bytes.NewBuffer(parsedBody))
	if err != nil {
		return Message{}, err
	}
	defer request.Body.Close()

	response, err := d.sendRequest(request)
	if err != nil {
		return Message{}, err
	}
	defer response.Body.Close()

	var jsonResponse Message
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return Message{}, err
	}

	if response.StatusCode != 200 {
		return Message{}, fmt.Errorf("%s (%d)", jsonResponse.Message, jsonResponse.Code)
	}

	return jsonResponse, nil
}
