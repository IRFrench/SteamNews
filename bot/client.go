package bot

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	ErrNoHello = errors.New("no hello message found")
)

type Bot struct {
	url          string
	conn         *websocket.Conn
	sequence     int
	requestChan  chan interface{}
	responseChan chan interface{}
}

func (b *Bot) Run(errChan chan<- error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	newConn, _, err := websocket.Dial(ctx, b.url, nil)
	if err != nil {
		errChan <- err
		return
	}

	b.conn = newConn

	for {
		select {
		case request := <-b.requestChan:
			if err := wsjson.Write(context.Background(), b.conn, request); err != nil {
				log.Err(err).Msg("failed to send message")
				errChan <- err
				b.Shutdown()
			}
		default:
			var response interface{}
			if err := wsjson.Read(context.Background(), b.conn, response); err != nil {
				log.Err(err).Msg("failed to read message")
				errChan <- err
				b.Shutdown()
			}
			b.responseChan <- response
		}
	}
}

func (b *Bot) Shutdown() {
	if b.conn != nil {
		b.conn.Close(websocket.StatusNormalClosure, "Shuting down service")
	}
}

func (b *Bot) Reconnect() {
	if err := b.conn.Close(websocket.StatusNoStatusRcvd, "Attempting reconnect"); err != nil {
		log.Err(err).Msg("failed to close the service")
		b.conn.CloseNow()
	}

}

func NewBot(socketUrl string) Bot {
	return Bot{
		url: socketUrl,
	}
}
