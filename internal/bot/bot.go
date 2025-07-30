package bot

import (
	"fmt"

	"copium-bot/internal/models"
)

type Opts struct {
	VoiceBot *Voice
}

type Bot struct {
	voice *Voice
}

func NewBot(opts Opts) *Bot {
	return &Bot{voice: opts.VoiceBot}
}

func (b *Bot) Process(r models.BotRequest) (models.BotResponse, error) {
	if !validRequest(r) {
		return models.BotResponse{}, fmt.Errorf("invalid request")
	}

	switch {
	case r.Message.Voice != nil:
		resp, err := b.voice.Process(r)
		if err != nil {
			return models.BotResponse{}, fmt.Errorf("voice: %w", err)
		}
		return resp, nil
	default:
		return models.BotResponse{}, nil
	}
}

func validRequest(r models.BotRequest) bool {
	switch {
	case r.User.ID == 0:
		return false
	case r.Message.ID == 0 || r.Message.ChatID == 0:
		return false
	default:
		return true
	}
}
