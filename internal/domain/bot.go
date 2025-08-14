package domain

type Processor interface {
	Process(request Request) (Response, error)
}

type Command interface {
	Processor
	Help() string
	ReactOn() []string
}

type User struct {
	ID   int64
	Name string
}

type Message struct {
	ID     int64
	ChatID int64

	Text      string
	Voice     []byte
	VideoNote []byte
	Command   string
}

type Request struct {
	User    User
	Message Message
}

type Response struct {
	ChatID  int64
	ReplyTo int64 // message id
	Text    string
}
