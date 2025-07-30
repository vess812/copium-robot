package main

import (
	"fmt"
	"log"

	"copium-bot/internal/bot"
	"copium-bot/internal/config"
	"copium-bot/internal/telegram"
	"copium-bot/internal/transcribe"

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

	transcriber, err := transcribe.NewTranscriber(transcribe.Opts{ModelPath: cfg.ModelPath})
	if err != nil {
		return fmt.Errorf("new transcriber: %w", err)
	}

	api, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	voice := bot.NewVoice(transcriber)
	b := bot.NewBot(bot.Opts{VoiceBot: voice})

	listener := telegram.NewListener(telegram.Opts{
		API: api,
		Bot: b,
	})

	log.Println("started")
	listener.Run()

	return nil
}
