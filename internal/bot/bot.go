package bot

import (
	"fmt"
	"strings"

	"copium-bot/internal/domain"
)

type Opts struct {
	Transcriber   domain.Processor
	CommandRouter domain.Processor
}

type Bot struct {
	transcriber   domain.Processor
	commandRouter domain.Processor
}

func NewBot(opts Opts) *Bot {
	return &Bot{
		transcriber:   opts.Transcriber,
		commandRouter: opts.CommandRouter,
	}
}

func (b *Bot) Process(r domain.Request) (domain.Response, error) {
	if !validRequest(r) {
		return domain.Response{}, fmt.Errorf("invalid request")
	}

	r.Message.Command = parseCommand(r)
	if r.Message.Command != "" {
		r.Message.Text = strings.TrimPrefix(r.Message.Text, fmt.Sprintf("%s ", r.Message.Command))
	}

	switch {
	case r.Message.Command != "":
		resp, err := b.commandRouter.Process(r)
		if err != nil {
			return domain.Response{}, fmt.Errorf("command: %w", err)
		}
		return resp, nil
	case r.Message.Voice != nil || r.Message.VideoNote != nil:
		resp, err := b.transcriber.Process(r)
		if err != nil {
			return domain.Response{}, fmt.Errorf("transcriber: %w", err)
		}
		return resp, nil
	default:
		return domain.Response{}, nil
	}
}

func validRequest(r domain.Request) bool {
	switch {
	case r.User.ID == 0:
		return false
	case r.Message.ID == 0 || r.Message.ChatID == 0:
		return false
	default:
		return true
	}
}

func parseCommand(r domain.Request) string {
	if !strings.HasPrefix(r.Message.Text, "!") {
		return ""
	}

	trim := strings.TrimPrefix(r.Message.Text, "!")
	split := strings.Split(trim, " ")
	if len(split) < 1 {
		return ""
	}

	return split[0]
}
