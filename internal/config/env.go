package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ModelPath        string `envconfig:"MODEL_PATH" required:"true"`
	TelegramBotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("process env: %w", err)
	}

	return cfg, nil
}
