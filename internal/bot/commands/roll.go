package commands

import (
	"math/rand/v2"
	"strconv"

	"copium-bot/internal/domain"
)

type Roll struct{}

func NewRoll() *Roll {
	return &Roll{}
}

func (r *Roll) Process(req domain.Request) (domain.Response, error) {
	return domain.Response{
		ChatID:  req.Message.ChatID,
		ReplyTo: req.Message.ID,
		Text:    strconv.Itoa(rand.IntN(101)),
	}, nil
}

func (r *Roll) Help() string {
	return "пишет случайное число в диапазоне от 0 до 100"
}

func (r *Roll) ReactOn() []string {
	return []string{"roll", "ролл"}
}
