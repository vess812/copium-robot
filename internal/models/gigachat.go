package models

import (
	"context"
	"errors"

	"github.com/Role1776/gigago"
)

type Gigachat struct {
	client *gigago.Client
	model  *gigago.GenerativeModel
}

const (
	systemInstruction = "Ты ассистент в групповом чате. Тебя зовут Копиум Дроид. Отвечай кратко и по делу."
	modelTemperature  = 1.0
)

func NewGigachat(client *gigago.Client, modelName string) *Gigachat {
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = systemInstruction
	model.Temperature = modelTemperature

	return &Gigachat{
		client: client,
		model:  model,
	}
}

func (g *Gigachat) Generate(ctx context.Context, text string) (string, error) {
	messages := []gigago.Message{
		{Role: gigago.RoleUser, Content: text},
	}

	resp, err := g.model.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) < 1 {
		return "", errors.New("no choices found")
	}

	return resp.Choices[0].Message.Content, nil
}
