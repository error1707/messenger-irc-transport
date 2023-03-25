package irc_client

import "time"

type (
	lseq uint64
	URI  string
)

type Message struct {
	MessageId lseq
	ParentId  lseq
	UserId    string
	Text      string
	Timestamp time.Time
}

type Transport interface {
	SendMessages(dest URI, msgs []Message) error
	ReceiveMessages(src URI, lastReceived lseq) ([]Message, error)
}

type ircTransport struct {
}

func (c *ircTransport) SendMessages(dest URI, msgs []Message) error {
	//TODO
	panic("implement me")
}

func (c *ircTransport) ReceiveMessages(src URI, lastReceived lseq) ([]Message, error) {
	//TODO
	panic("implement me")
}
