package irc_transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"go.uber.org/multierr"
	"gopkg.in/irc.v4"
)

type Message struct {
	MessageId uint64    `json:"message_id"`
	ParentId  uint64    `json:"parent_id"`
	UserId    string    `json:"user_id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type Transport interface {
	SendMessages(dest string, msgs []Message) error
	ReceiveMessages(src string, lastReceived uint64) ([]Message, error)
}

const MessageCommand = "PRIVMSG"

type ircTransport struct {
	ircClient     *irc.Client
	defaultPrefix *irc.Prefix
}

func NewIRCTransport(username string) (Transport, error) {
	conn, err := net.Dial("tcp", "bitcoin.uk.eu.dal.net:6667")
	if err != nil {
		return nil, fmt.Errorf("can't connect to server: %w", err)
	}
	transport := &ircTransport{
		defaultPrefix: &irc.Prefix{
			Name: username,
			User: username,
			Host: username,
		},
	}
	transport.ircClient = irc.NewClient(conn, irc.ClientConfig{
		Nick:          username,
		User:          username,
		Name:          username,
		PingFrequency: 10 * time.Second,
		PingTimeout:   5 * time.Second,
		Handler:       irc.HandlerFunc(transport.messageHandler),
	})
	transport.ircClient.CapRequest("message-tags", true)
	transport.ircClient.Conn.Writer.DebugCallback = func(line string) {
		log.Printf("[DEBUG] %s\n", line)
	}
	//transport.ircClient.Conn.Reader.DebugCallback = func(line string) {
	//	log.Printf("[DEBUG] %s\n", line)
	//}
	go func() {
		err := transport.ircClient.Run()
		if err != nil {
			log.Printf("irc transport run finished with error: %v", err)
		}
	}()
	return transport, nil
}

func (t *ircTransport) SendMessages(dest string, msgs []Message) error {
	var errs error
	for _, msg := range msgs {
		body, _ := json.Marshal(msg)
		err := t.ircClient.WriteMessage(&irc.Message{
			Prefix:  t.defaultPrefix.Copy(),
			Command: MessageCommand,
			Params: []string{
				dest, string(body),
			},
		})
		errs = multierr.Append(errs, err)
	}
	return errs
}

func (t *ircTransport) ReceiveMessages(src string, lastReceived uint64) ([]Message, error) {
	//TODO
	panic("implement me")
}

func (t *ircTransport) messageHandler(client *irc.Client, msg *irc.Message) {
	log.Printf("[RECEIVED] %s\n", msg.String())
}
