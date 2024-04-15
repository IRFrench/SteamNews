package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendMessage struct {
	Content string `json:"content,omitempty"`
}

type ReactionEmoji struct {
	ID   *string `json:"id"`
	Name string  `json:"name"`
}

type ReactionCountDetails struct {
	Burst  int `json:"burst"`
	Normal int `json:"normal"`
}

type Reaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details"`
	Me           bool                 `json:"me"`
	Emoji        ReactionEmoji        `json:"emoji"`
}

type Author struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
	Avatar        string `json:"avatar"`
}

type MessageObject struct {
	Reactions       []Reaction `json:"reactions,omitempty"`
	TTS             bool       `json:"tts,omitempty"`
	Timestamp       string     `json:"timestamp,omitempty"`
	MentionEveryone bool       `json:"mention_everyone,omitempty"`
	ID              string     `json:"id,omitempty"`
	Pinned          bool       `json:"pinned,omitempty"`
	EditedTimestamp *string    `json:"edited_timestamp,omitempty"`
	Author          Author     `json:"author,omitempty"`
	Content         string     `json:"content,omitempty"`
	ChannelID       string     `json:"channel_id,omitempty"`
	Type            int        `json:"type,omitempty"`
	Message         string     `json:"message,omitempty"`
	Code            int        `json:"code,omitempty"`
}

func (d *DiscordClient) SendNewsMessage(news []Game, channelId string) error {
	discordUrl := fmt.Sprintf("%s/v%d/channels/%s/messages", discordUrl, discordApiVersion, channelId)

	messages := formatMessage(news)
	for _, message := range messages {

		parsedBody, err := json.Marshal(SendMessage{Content: message})
		if err != nil {
			return err
		}

		request, err := http.NewRequest(http.MethodPost, discordUrl, bytes.NewBuffer(parsedBody))
		if err != nil {
			return err
		}
		defer request.Body.Close()

		response, err := d.sendRequest(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		var jsonResponse MessageObject
		if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
			return err
		}

		if response.StatusCode != 200 {
			return fmt.Errorf("%s (%d)", jsonResponse.Message, jsonResponse.Code)
		}
	}

	return nil
}
