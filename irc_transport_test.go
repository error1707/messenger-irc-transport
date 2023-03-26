package irc_transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IRCTransportTestSuite struct {
	suite.Suite

	client Transport
}

func TestIRCTransportTestSuite(t *testing.T) {
	suite.Run(t, new(IRCTransportTestSuite))
}

func (s *IRCTransportTestSuite) SetupSuite() {
	var err error
	s.client, err = NewIRCTransport("test_user")
	s.Require().NoError(err)
}

func (s *IRCTransportTestSuite) TestSendMessages() {
	time.Sleep(5 * time.Second)
	err := s.client.SendMessages("error1707_", []Message{
		{
			MessageId: 1,
			ParentId:  1,
			UserId:    "test_user",
			Text:      "some text",
			Timestamp: time.Now(),
		},
	})
	s.Require().NoError(err)
	time.Sleep(20 * time.Second)
}
