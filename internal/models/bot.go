package models

type Bot interface {
	Process(request BotRequest) (BotResponse, error)
}

type User struct {
	ID   int64
	Name string
}

type Message struct {
	ID     int64
	ChatID int64

	Text  string
	Voice []byte
}

type BotRequest struct {
	User    User
	Message Message
}

type BotResponse struct {
	ChatID  int64
	ReplyTo int64 // message id
	Text    string
}
