package irc_transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IRCTransportTestSuite struct {
	suite.Suite

	client1 Transport
	client2 Transport
}

func TestIRCTransportTestSuite(t *testing.T) {
	suite.Run(t, new(IRCTransportTestSuite))
}

func (s *IRCTransportTestSuite) SetupSuite() {
	var err error
	s.client1, err = NewIRCTransport("test_user_1")
	s.Require().NoError(err)
	s.client2, err = NewIRCTransport("test_user_2")
	s.Require().NoError(err)
}

func (s *IRCTransportTestSuite) TestSendMessages() {
	msgs, err := s.client2.ReceiveMessages("test_user_1", 0)
	s.Require().NoError(err)
	s.Assert().Len(msgs, 0)

	err = s.client1.SendMessages("test_user_2", []Message{
		{
			MessageId: 1,
			ParentId:  1,
			UserId:    "test_user_2",
			Text:      "some text",
			Timestamp: time.Now(),
		},
	})
	s.Require().NoError(err)
	time.Sleep(2 * time.Second)

	msgs, err = s.client2.ReceiveMessages("test_user_1", 0)
	s.Require().NoError(err)
	s.Assert().Len(msgs, 1)
}
