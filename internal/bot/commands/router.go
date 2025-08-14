package commands

import (
	"fmt"
	"slices"
	"strings"

	"copium-bot/internal/domain"
)

type Router struct {
	commands []domain.Command
}

func NewRouter(commands ...domain.Command) *Router {
	return &Router{
		commands: commands,
	}
}

func (r *Router) Process(req domain.Request) (domain.Response, error) {
	switch req.Message.Command {
	case "":
		return domain.Response{}, fmt.Errorf("missing command")
	case "help", "помощь", "хелп":
		return r.buildHelpResponse(req)
	default:
		return r.processCommand(req)
	}
}

func (r *Router) buildHelpResponse(req domain.Request) (domain.Response, error) {
	b := strings.Builder{}
	b.WriteString("Доступные команды:\n")
	b.WriteString("!help, !помощь, !хелп: вывести это сообщение\n")
	for _, c := range r.commands {
		b.WriteString(buildHelpLine(c.ReactOn(), c.Help()))
	}
	return domain.Response{
		ChatID:  req.Message.ChatID,
		ReplyTo: req.Message.ID,
		Text:    b.String(),
	}, nil
}

func buildHelpLine(commands []string, help string) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("!%s", commands[0]))
	for _, c := range commands[1:] {
		b.WriteString(fmt.Sprintf(", !%s", c))
	}
	b.WriteString(fmt.Sprintf(": %s\n", help))
	return b.String()
}

func (r *Router) processCommand(req domain.Request) (domain.Response, error) {
	for _, c := range r.commands {
		if slices.Contains(c.ReactOn(), req.Message.Command) {
			return c.Process(req)
		}
	}

	return domain.Response{}, fmt.Errorf("command not found")
}
