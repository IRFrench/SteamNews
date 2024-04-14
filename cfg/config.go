package cfg

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Steam   SteamConfig   `yaml:"steam,omitempty"`
	Discord DiscordConfig `yaml:"discord,omitempty"`
}

type SteamConfig struct {
	User    User  `yaml:"user,omitempty"`
	Added   []int `yaml:"added,omitempty"`
	Removed []int `yaml:"removed,omitempty"`
}

type User struct {
	Key string `yaml:"key,omitempty"`
	Id  int    `yaml:"id,omitempty"`
}

type DiscordConfig struct {
	BotToken string `yaml:"bot_token,omitempty"`
	UserId   int    `yaml:"user_id,omitempty"`
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
