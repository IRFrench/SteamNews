package main

import (
	"SteamNews/bot/discord"
	"SteamNews/bot/steam"
	"SteamNews/cfg"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	configEnviroment = "ENV"
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

	log.Info().Msg("starting the service")

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

	ticker := time.NewTicker(1 * time.Second)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config, err := cfg.ReadConfiguration(configPath)
	if err != nil {
		log.Err(err).
			Str("config path", configPath).
			Msg("failed to load configuration")
		return err
	}

	steamUser := steam.NewUser(
		config.Steam.User.Key,
		config.Steam.User.Id,
	)
	steamClient := steam.NewClient(steamUser)
	log.Info().
		Str("key", config.Steam.User.Key).
		Int("id", config.Steam.User.Id).
		Msg("created steam client")

	discordClient := discord.NewDiscordClient(config.Discord.BotToken, config.Discord.UserId)
	log.Info().
		Str("token", config.Discord.BotToken).
		Int("user id", config.Discord.UserId).
		Msg("created discord client")

	for {
		select {
		case signal := <-sigChan:
			log.Info().Str("signal", signal.String()).Msg("recieved end signal")
			return nil
		case <-ticker.C:
			if err := SendNewsUpdate(&steamClient, &discordClient); err != nil {
				log.Err(err).Msg("failed to send news update")
			}
			return nil
		}
	}
}

func SendNewsUpdate(steamClient *steam.SteamClient, discordClient *discord.DiscordClient) error {
	lastTime := time.Now().Add(-24 * time.Hour)

	channel, err := discordClient.CreateDmChannel()
	if err != nil {
		return err
	}
	fmt.Println("response")
	fmt.Println(channel.Id, channel.Type)
	return fmt.Errorf("test")

	log.Debug().Msg("collecting owned games")
	games, err := steamClient.GetOwnedGames()
	if err != nil {
		return err
	}
	log.Debug().Int("game count", len(games)).Msg("collected owned games")

	collectedNews := []discord.Game{}
	for _, game := range games {
		log.Debug().
			Str("game", game.Name).
			Int("id", game.Appid).
			Msg("collecting game news")

		newsArticles, err := steamClient.GetAppNews(game.Appid)
		if err != nil {
			return err
		}

		log.Debug().Msg("collected games news")

		var articles []discord.ShortArticle
		for _, article := range newsArticles {
			articleDate := time.Unix(int64(article.Date), 0)
			log.Debug().
				Str("title", article.Title).
				Time("date", articleDate).
				Bool("new", articleDate.After(lastTime)).
				Msg("reading article")

			if articleDate.After(lastTime) {
				articles = append(articles, discord.ShortArticle{
					Title:    article.Title,
					Url:      article.Url,
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

		if len(articles) > 0 {
			collectedNews = append(collectedNews, discord.Game{
				Name: game.Name,
				News: articles,
			})
		}
	}

	log.Debug().
		Int("game news count", len(collectedNews)).
		Msg("Collected all news for games")

	return nil
}
