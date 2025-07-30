package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"copium-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Opts struct {
	API *tgbotapi.BotAPI
	Bot models.Bot
}

type Listener struct {
	api *tgbotapi.BotAPI
	bot models.Bot
}

func NewListener(opts Opts) *Listener {
	return &Listener{
		api: opts.API,
		bot: opts.Bot,
	}
}

func (l *Listener) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := l.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

			req, err := l.parseUpdate(update)
			if err != nil {
				log.Println(fmt.Sprintf("parse update: %v", err))
			}

			resp, err := l.bot.Process(req)
			if err != nil {
				log.Println(fmt.Sprintf("process: %v", err))
			}

			err = l.sendResponse(resp)
			if err != nil {
				log.Println(fmt.Sprintf("send response: %v", err))
			}
		}
	}
}

func (l *Listener) parseUpdate(u tgbotapi.Update) (models.BotRequest, error) {
	r := models.BotRequest{
		User: models.User{
			ID:   u.Message.From.ID,
			Name: u.Message.From.UserName,
		},
		Message: models.Message{
			ID:     int64(u.Message.MessageID),
			ChatID: u.Message.Chat.ID,
			Text:   u.Message.Text,
		},
	}

	if u.Message.Voice != nil {
		f, err := l.api.GetFile(tgbotapi.FileConfig{FileID: u.Message.Voice.FileID})
		if err != nil {
			return models.BotRequest{}, fmt.Errorf("get file: %w", err)
		}
		voice, err := downloadVoice(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", l.api.Token, f.FilePath))
		if err != nil {
			return models.BotRequest{}, fmt.Errorf("download voice: %w", err)
		}
		r.Message.Voice = voice
	}

	return r, nil
}

func downloadVoice(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}
	return buf.Bytes(), nil
}

func (l *Listener) sendResponse(r models.BotResponse) error {
	if r.ChatID == 0 {
		return nil
	}

	if len(r.Text) == 0 {
		return errors.New("empty response")
	}

	msg := tgbotapi.NewMessage(r.ChatID, r.Text)

	if r.ReplyTo != 0 {
		msg.ReplyToMessageID = int(r.ReplyTo)
	}

	if _, err := l.api.Send(msg); err != nil {
		return fmt.Errorf("api: %w", err)
	}

	return nil
}
