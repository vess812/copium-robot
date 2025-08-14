package commands

import (
	"context"
	"fmt"

	"copium-bot/internal/domain"
)

type Generator interface {
	Generate(ctx context.Context, text string) (string, error)
}

type AI struct {
	generator Generator
}

func NewAI(generator Generator) *AI {
	return &AI{generator: generator}
}

func (a *AI) Process(r domain.Request) (domain.Response, error) {
	resp, err := a.generator.Generate(context.Background(), r.Message.Text)
	if err != nil {
		return domain.Response{}, fmt.Errorf("generate: %w", err)
	}

	return domain.Response{
		ChatID:  r.Message.ChatID,
		ReplyTo: r.Message.ID,
		Text:    resp,
	}, nil
}

func (a *AI) Help() string {
	return "отправить промпт в нейронку"
}

func (a *AI) ReactOn() []string {
	return []string{"ai", "аи"}
}
