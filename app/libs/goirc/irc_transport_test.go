package irc_transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IRCTransportTestSuite struct {
	suite.Suite

	client1 *IrcTransport
	client2 *IrcTransport
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
	someMessage := "some message"
	waiter := make(chan struct{})
	go func() {
		s.client2.StartReceiveMessagesFrom("test_user_1")
		<-time.After(time.Second)
		msg, err := s.client2.GetMessageFrom("test_user_1")
		s.Require().NoError(err)
		s.Assert().Equal(someMessage, msg)
		close(waiter)
	}()

	err := s.client1.SendMessages("test_user_2", someMessage)
	s.Require().NoError(err)
	select {
	case <-waiter:
	case <-time.After(5 * time.Second):
		s.FailNow("wait for message to long")
	}
}
