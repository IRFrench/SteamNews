package discord

import (
	"fmt"
	"strings"
	"time"
)

type Game struct {
	Name string
	Id   int
	News []ShortArticle
}

type ShortArticle struct {
	Title    string
	Url      string
	Author   string
	Contents string
	Date     time.Time
}

func formatMessage(games []Game) []string {
	messageList := []string{}

	var message strings.Builder
	for _, game := range games {

		var tempMessage strings.Builder
		tempMessage.WriteString(fmt.Sprintf("**%s (%d)**\n", game.Name, game.Id))

		for _, article := range game.News {
			tempMessage.WriteString(fmt.Sprintf(
				"%s\n",
				article.Url,
			))
		}

		// Discord has a max length of 2000 characters in a message
		if (message.Len() + tempMessage.Len()) > 2000 {
			messageList = append(messageList, message.String())
			message.Reset()
		}
		message.WriteString(tempMessage.String())
	}
	messageList = append(messageList, message.String())

	return messageList
}
