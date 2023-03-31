package irc_transport

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

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

func (t *IrcTransport) SendMessages(dest string, msg string) error {
	return t.ircClient.WriteMessage(&irc.Message{
		Prefix:  t.defaultPrefix.Copy(),
		Command: MessageCommand,
		Params: []string{
			dest, msg,
		},
	})
}

func (t *IrcTransport) StartReceiveMessagesFrom(src string) {
	t.receiveChannels.Store(src, make(chan string, MessageReceiveBufferSize))
}

func (t *IrcTransport) StopReceiveMessagesFrom(src string) {
	t.receiveChannels.Delete(src)
}

func (t *IrcTransport) GetMessageFrom(src string) (string, error) {
	msgChan, ok := t.receiveChannels.Load(src)
	if !ok {
		return "", errors.New("source is not listened")
	}
	select {
	case msg := <-msgChan.(chan string):
		return msg, nil
	default:
		return "", io.EOF
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
