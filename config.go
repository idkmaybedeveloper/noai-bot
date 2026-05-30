package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`

	GitHub struct {
		AppID         int64  `toml:"app_id"`
		PrivateKey    string `toml:"private_key_path"`
		WebhookSecret string `toml:"webhook_secret"`
	} `toml:"github"`

	ExtraPatterns  []string `toml:"extra_patterns"`
	DeclineMessage string   `toml:"decline_message"`
}

func loadConfig(path string) (*Config, error) {
	cfg := Config{}
	cfg.Server.Port = 8080
	cfg.DeclineMessage = defaultDeclineMessage

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

const defaultDeclineMessage = `## noai bot

this pull request has been automatically closed because one or more commits were co-authored by an ai assistant

**policy:** this project does not accept ai-generated contributions

if you believe this is a mistake, please contact a maintainer`
