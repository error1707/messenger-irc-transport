package irc_transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"go.uber.org/multierr"
	"gopkg.in/irc.v4"
)

type Message struct {
	MessageId int64    `json:"message_id"`
	ParentId  int64    `json:"parent_id"`
	UserId    string    `json:"user_id"`
	Text      string    `json:"text"`
	Timestamp int64 `json:"timestamp"`
}

const MessageCommand = "PRIVMSG"
const MessageReceiveBufferSize = 32

type IrcTransport struct {
	ircClient       *irc.Client
	defaultPrefix   *irc.Prefix
	receiveChannels sync.Map

	initialized chan struct{}
}

func NewIrcTransport(username string) (*IrcTransport, error) {
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
	//transport.ircClient.Conn.Reader.DebugCallback = func(line string) {
	//	log.Printf("[DEBUG] %s\n", line)
	//}
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

func (t *IrcTransport) SendMessages(dest string, msgs *Message) error {
	var errs error
	var msg = *msgs
	{
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

func (t *IrcTransport) ReceiveMessages(src string, lastReceived uint64) (Message, error) {
	rawValue, _ := t.receiveChannels.LoadOrStore(src, make(chan *Message, MessageReceiveBufferSize))
	msgChan := rawValue.(chan *Message)
	var received []Message
outer:
	for {
		select {
		case msg := <-msgChan:
			received = append(received, *msg)
		default:
			break outer
		}
	}
	return received[0], nil
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
	msgChan := rawValue.(chan *Message)
	clientMsg := &Message{}
	err := json.Unmarshal([]byte(msg.Trailing()), clientMsg)
	if err != nil {
		log.Printf("[ERROR] Can't parse message: %v\n", err)
	}
	log.Printf("[INFO: %s] Writing to channel: %s\n", t.defaultPrefix.User, msg.Name)
	msgChan <- clientMsg
}
