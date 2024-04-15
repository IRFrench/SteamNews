package cfg

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Steam   SteamConfig   `yaml:"steam,omitempty"`
	Discord DiscordConfig `yaml:"discord,omitempty"`
	Users   []User        `yaml:"users,omitempty"`
}

type SteamConfig struct {
	Key string `yaml:"key,omitempty"`
}

type DiscordConfig struct {
	BotToken string `yaml:"bot_token,omitempty"`
}

type User struct {
	Name      string    `yamk:"name,omitempty"`
	DiscordId int       `yaml:"discord_id,omitempty"`
	Steam     SteamUser `yaml:"steam,omitempty"`
}

type SteamUser struct {
	Id      int         `yaml:"id,omitempty"`
	Added   []SteamGame `yaml:"added,omitempty"`
	Removed []int       `yaml:"removed,omitempty"`
}

type SteamGame struct {
	Name string `yaml:"name,omitempty"`
	Id   int    `yaml:"id,omitempty"`
}

func ReadConfiguration(configPath string) (Configuration, error) {
	conf, err := os.ReadFile(configPath)
	if err != nil {
		return Configuration{}, err
	}

	var parsedConf Configuration
	if err := yaml.Unmarshal(conf, &parsedConf); err != nil {
		return Configuration{}, err
	}

	return parsedConf, nil
}
