package bot

import (
	"fmt"

	"copium-bot/internal/domain"
)

type Opts struct {
	Transcriber domain.Bot
}

type Bot struct {
	transcriber domain.Bot
}

func NewBot(opts Opts) *Bot {
	return &Bot{transcriber: opts.Transcriber}
}

func (b *Bot) Process(r domain.BotRequest) (domain.BotResponse, error) {
	if !validRequest(r) {
		return domain.BotResponse{}, fmt.Errorf("invalid request")
	}

	switch {
	case r.Message.Voice != nil || r.Message.VideoNote != nil:
		resp, err := b.transcriber.Process(r)
		if err != nil {
			return domain.BotResponse{}, fmt.Errorf("transcriber: %w", err)
		}
		return resp, nil
	default:
		return domain.BotResponse{}, nil
	}
}

func validRequest(r domain.BotRequest) bool {
	switch {
	case r.User.ID == 0:
		return false
	case r.Message.ID == 0 || r.Message.ChatID == 0:
		return false
	default:
		return true
	}
}
