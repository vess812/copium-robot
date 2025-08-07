package main

import (
	"fmt"
	"log"

	"copium-bot/internal/bot"
	"copium-bot/internal/config"
	"copium-bot/internal/model"
	"copium-bot/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := app()
	if err != nil {
		log.Fatal(fmt.Errorf("app: %w", err))
	}
}

func app() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	voskModel, err := model.NewVoskModel(model.Opts{ModelPath: cfg.ModelPath})
	if err != nil {
		return fmt.Errorf("new vosk model: %w", err)
	}

	api, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	transcriber := bot.NewTranscriber(voskModel)
	b := bot.NewBot(bot.Opts{Transcriber: transcriber})

	listener := telegram.NewListener(telegram.Opts{
		API: api,
		Bot: b,
	})

	log.Println("started")
	listener.Run()

	return nil
}
