package bot

import (
	"math/rand/v2"
	"time"

	"github.com/rs/zerolog/log"
)

func (b *Bot) performHeartbeat(interval int) {
	jitter := rand.IntN(1000)
	time.Sleep(time.Duration(jitter) * time.Microsecond)

	newTicker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	heartbeat := HeartbeatMessage{
		EventPayload: EventPayload{
			Op: 1,
		},
	}

	for {
		<-newTicker.C
		if b.sequence != 0 {
			heartbeat.D = b.sequence
		}
		log.Debug().
			Int("sequence", heartbeat.D).
			Msg("sending heartbeat")

		b.requestChan <- heartbeat

		response := <-b.responseChan
		ack, ok := response.(EventPayload)

		if !ok {
			log.Error().Msg("could not convert event payload")
			continue
		}

		if ack.Op == 11 {
			continue
		}
	}
}
