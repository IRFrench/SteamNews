package main

import (
	"SteamNews/bot/discord"
	"SteamNews/bot/steam"
	"SteamNews/cfg"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	configEnviroment = "ETC"
)

func main() {
	if err := runServer(); err != nil {
		log.Err(err).Msg("unexpected error")
		os.Exit(1)
	}
	os.Exit(0)
}

func runServer() error {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("creating the service")

	flags := cfg.ReadFlags()
	if flags.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug turned on")
	}

	configPath, ok := os.LookupEnv(configEnviroment)
	if !ok {
		log.Info().
			Str("environment", configEnviroment).
			Msg("missing environment variable")
		return fmt.Errorf("missing %s environment", configEnviroment)
	}

	config, err := cfg.ReadConfiguration(configPath)
	if err != nil {
		log.Err(err).
			Str("config path", configPath).
			Msg("failed to load configuration")
		return err
	}

	if !flags.Quick {
		waitForStart := time.NewTicker(2 * time.Second)
		log.Info().
			Int("hour", config.StartTime.Hour).
			Int("minute", config.StartTime.Minute).
			Msg("Waiting to start the service")
		for {
			<-waitForStart.C
			currentHour, currentMinute, _ := time.Now().Clock()
			if currentHour == config.StartTime.Hour && currentMinute == config.StartTime.Minute {
				break
			}
		}
	}

	log.Info().Time("start", time.Now()).Msg("starting the service")

	ticker := time.NewTicker(24 * time.Hour)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	steamClient := steam.NewClient(config.Steam.Key)
	log.Info().
		Str("key", config.Steam.Key).
		Msg("created steam client")

	discordClient := discord.NewDiscordClient(config.Discord.BotToken)
	log.Info().
		Str("token", config.Discord.BotToken).
		Msg("created discord client")

	// Test the service to make sure it works.
	// If the service fails on the first run, it likely won't work again so we
	// error out at that point.
	log.Info().Msg("testing service")
	if err := SendNewsUpdate(&steamClient, &discordClient, config.Users); err != nil {
		log.Err(err).Msg("failed to send news update")
		return err
	}
	log.Info().Msg("test complete")

	for {
		select {
		case signal := <-sigChan:
			log.Info().Str("signal", signal.String()).Msg("recieved end signal")
			return nil
		case <-ticker.C:
			if err := SendNewsUpdate(&steamClient, &discordClient, config.Users); err != nil {
				log.Err(err).Msg("failed to send news update")
			}
		}
	}
}

func SendNewsUpdate(steamClient *steam.SteamClient, discordClient *discord.DiscordClient, users []cfg.User) error {
	lastTime := time.Now().Add(-24 * time.Hour)
	log.Info().Msg("Gathering news for users")

	for _, user := range users {
		log.Info().Str("name", user.Name).Msg("collecting news")

		// Collect news from Steam
		log.Debug().Msg("collecting owned games")
		games, err := steamClient.GetOwnedGames(user.Steam.Id)
		if err != nil {
			return err
		}
		log.Debug().Int("game count", len(games)).Msg("collected owned games")

		collectedNews := []discord.Game{}

		// Add added games to the game list
		for _, addedGame := range user.Steam.Added {
			log.Debug().Str("game", addedGame.Name).Msg("added game")
			games = append(games, steam.Game{
				Name:  addedGame.Name,
				Appid: addedGame.Id,
			})
		}

		for _, game := range games {
			// If it has been removed, move on
			if contains(user.Steam.Removed, game.Appid) {
				log.Warn().Str("game", game.Name).Msg("skipped game")
				continue
			}

			log.Debug().
				Msg("collecting game news")

			newsArticles, err := steamClient.GetAppNews(game.Appid)
			if err != nil {
				return err
			}

			log.Debug().
				Str("game", game.Name).
				Int("id", game.Appid).
				Msg("collected games news")

			// Sort through the news
			var articles []discord.ShortArticle
			for _, article := range newsArticles {
				articleDate := time.Unix(int64(article.Date), 0)

				if articleDate.After(lastTime) {
					url, err := url.Parse(article.Url)
					if err != nil {
						log.Warn().Str("url", article.Url).Msg("could not parse url")
					}

					articles = append(articles, discord.ShortArticle{
						Title:    article.Title,
						Url:      url.String(),
						Author:   article.Author,
						Contents: article.Contents,
						Date:     articleDate,
					})
				}
			}
			log.Debug().
				Int("articles", len(articles)).
				Str("game", game.Name).
				Msg("collected articles")

			// If there is no new news, skip
			if len(articles) > 0 {
				collectedNews = append(collectedNews, discord.Game{
					Name: game.Name,
					Id:   game.Appid,
					News: articles,
				})
			}
		}

		log.Debug().
			Int("game news count", len(collectedNews)).
			Msg("collected all news for games")

		// Send discord message
		channel, err := discordClient.CreateDmChannel(user.DiscordId)
		if err != nil {
			return err
		}
		log.Debug().
			Str("channel", channel.Id).
			Msg("created dm channel")

		for _, recipient := range channel.Recipients {
			log.Info().
				Str("username", recipient.Username).
				Msg("sending message")
		}

		if err := discordClient.SendNewsMessage(collectedNews, channel.Id); err != nil {
			return err
		}

	}

	log.Info().
		Msg("sent messages to users")

	return nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
