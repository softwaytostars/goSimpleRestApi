package emails

import (
	"errors"
	"sync"
	"time"
)

/*
	simple implementation for tests
*/
type SimpleSmtpConnectorImpl struct {
	mutex           sync.Mutex
	ErrorConnect    bool
	ErrorDisconnect bool
	IsConnected     bool
	NSent           int
}

func (s *SimpleSmtpConnectorImpl) ConnectionIsOpen() bool {
	return s.IsConnected
}

func (s *SimpleSmtpConnectorImpl) Connect() error {
	if s.ErrorConnect {
		return errors.New("cannot connect")
	}
	s.IsConnected = true
	return nil
}

func (s *SimpleSmtpConnectorImpl) Disconnect() error {
	if s.ErrorDisconnect {
		return errors.New("cannot connect")
	}
	s.IsConnected = false
	return nil
}

func (s *SimpleSmtpConnectorImpl) Send(m *EmailMessage) error {
	//use mutex here in order to use RLock instead of Lock in EmailSender::Send . Indeed, the send method should not change the state.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	time.Sleep(30 * time.Millisecond)
	if !s.IsConnected {
		return errors.New("cannot send")
	}

	s.NSent++
	return nil
}
