package main

import (
	"context"
	"fmt"
	"log"

	"copium-bot/internal/bot"
	"copium-bot/internal/bot/commands"
	"copium-bot/internal/config"
	"copium-bot/internal/models"
	"copium-bot/internal/telegram"

	"github.com/Role1776/gigago"
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

	api, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("new bot api: %w", err)
	}

	voskModel, err := models.NewVoskModel(models.Opts{ModelPath: cfg.ModelPath})
	if err != nil {
		return fmt.Errorf("new vosk model: %w", err)
	}
	transcriber := bot.NewTranscriber(voskModel)

	client, err := gigago.NewClient(context.Background(), cfg.GigachatAPIKey, gigago.WithCustomInsecureSkipVerify(true))
	if err != nil {
		return fmt.Errorf("new gigachat client: %w", err)
	}
	defer client.Close()
	gigachat := models.NewGigachat(client, "GigaChat")
	commandRouter := commands.NewRouter(commands.NewRoll(), commands.NewAI(gigachat))

	b := bot.NewBot(bot.Opts{
		Transcriber:   transcriber,
		CommandRouter: commandRouter,
	})

	listener := telegram.NewListener(telegram.Opts{
		API: api,
		Bot: b,
	})

	log.Println("started")
	listener.Run()

	return nil
}
