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
	errorConnect    bool
	errorDisconnect bool
	isConnected     bool
	nSent           int
}

func (s *SimpleSmtpConnectorImpl) ConnectionIsOpen() bool {
	return s.isConnected
}

func (s *SimpleSmtpConnectorImpl) Connect() error {
	if s.errorConnect {
		return errors.New("Cannot connect")
	}
	s.isConnected = true
	return nil
}

func (s *SimpleSmtpConnectorImpl) Disconnect() error {
	if s.errorDisconnect {
		return errors.New("Cannot connect")
	}
	s.isConnected = false
	return nil
}

func (s *SimpleSmtpConnectorImpl) Send(m *EmailMessage) error {
	//use mutex here in order to use RLock instead of Lock in EmailSender::Send . Indeed, the send method should not change the state.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	time.Sleep(30 * time.Millisecond)
	if !s.isConnected {
		return errors.New("Cannot send")
	}

	s.nSent++
	return nil
}
