package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"copium-bot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Opts struct {
	API *tgbotapi.BotAPI
	Bot domain.Processor
}

type Listener struct {
	api *tgbotapi.BotAPI
	bot domain.Processor
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
				log.Printf("parse update: %v\n", err)
			}

			resp, err := l.bot.Process(req)
			if err != nil {
				log.Printf("process: %v\n", err)
			}

			err = l.sendResponse(resp)
			if err != nil {
				log.Printf("send response: %v\n", err)
			}
		}
	}
}

func (l *Listener) parseUpdate(u tgbotapi.Update) (domain.Request, error) {
	r := domain.Request{
		User: domain.User{
			ID:   u.Message.From.ID,
			Name: u.Message.From.UserName,
		},
		Message: domain.Message{
			ID:     int64(u.Message.MessageID),
			ChatID: u.Message.Chat.ID,
			Text:   u.Message.Text,
		},
	}

	if u.Message.Voice != nil || u.Message.VideoNote != nil {
		switch {
		case u.Message.Voice != nil:
			data, err := l.downloadTelegramFile(u.Message.Voice.FileID)
			if err != nil {
				return r, err
			}
			r.Message.Voice = data
		case u.Message.VideoNote != nil:
			data, err := l.downloadTelegramFile(u.Message.VideoNote.FileID)
			if err != nil {
				return r, err
			}
			r.Message.VideoNote = data
		}
	}

	return r, nil
}

func (l *Listener) downloadTelegramFile(fileID string) ([]byte, error) {
	f, err := l.api.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}
	data, err := downloadFile(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", l.api.Token, f.FilePath))
	if err != nil {
		return nil, fmt.Errorf("download voice: %w", err)
	}

	return data, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}
	return buf.Bytes(), nil
}

func (l *Listener) sendResponse(r domain.Response) error {
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
