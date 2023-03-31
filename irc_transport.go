package irc_transport

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"go.uber.org/multierr"
	"gopkg.in/irc.v4"
)

const MessageCommand = "PRIVMSG"
const MessageReceiveBufferSize = 32

type IrcTransport struct {
	ircClient       *irc.Client
	defaultPrefix   *irc.Prefix
	receiveChannels sync.Map

	initialized chan struct{}
}

func NewIRCTransport(username string) (*IrcTransport, error) {
	conn, err := net.Dial("tcp", "bitcoin.uk.eu.dal.net:6667")
	if err != nil {
		return nil, fmt.Errorf("can't connect to server: %w", err)
	}
	transport := &IrcTransport{
		defaultPrefix: &irc.Prefix{
			Name: username,
			User: username,
			Host: username,
		},
		initialized: make(chan struct{}),
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
		log.Printf("[SENT: %s] %s\n", username, line)
	}
	go func() {
		err := transport.ircClient.Run()
		if err != nil {
			log.Printf("irc transport run finished with error: %v", err)
		}
	}()
	select {
	case <-transport.initialized:
	case <-time.After(5 * time.Second):
		return nil, errors.New("client initializing to long")
	}
	return transport, nil
}

func (t *IrcTransport) SendMessages(dest string, msgs []string) error {
	var errs error
	for _, msg := range msgs {
		err := t.ircClient.WriteMessage(&irc.Message{
			Prefix:  t.defaultPrefix.Copy(),
			Command: MessageCommand,
			Params: []string{
				dest, msg,
			},
		})
		errs = multierr.Append(errs, err)
	}
	return errs
}

func (t *IrcTransport) ReceiveMessages(src string, handler func(string) bool) {
	msgChan := make(chan string, MessageReceiveBufferSize)
	t.receiveChannels.Store(src, msgChan)
	defer t.receiveChannels.Delete(src)
	for msg := range msgChan {
		if !handler(msg) {
			break
		}
	}
}

func (t *IrcTransport) messageHandler(client *irc.Client, msg *irc.Message) {
	log.Printf("[RECEIVED: %s] %s\n", t.defaultPrefix.User, msg.String())
	if msg.Command == "MODE" {
		select {
		case <-t.initialized:
		default:
			close(t.initialized)
		}
	}
	if msg.Command != MessageCommand {
		return
	}
	rawValue, ok := t.receiveChannels.Load(msg.Name)
	if !ok {
		log.Printf("[WARN: %s] Client not listening for messages from: %s\n", t.defaultPrefix.User, msg.Name)
		return
	}
	log.Printf("[INFO: %s] Writing to channel: %s\n", t.defaultPrefix.User, msg.Name)
	rawValue.(chan string) <- msg.Trailing()
}
