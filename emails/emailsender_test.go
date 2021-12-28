package emails

import (
	"fmt"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type SmtpConnectorMock struct {
	mock.Mock
}

func (s *SmtpConnectorMock) ConnectionIsOpen() bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *SmtpConnectorMock) Connect() error {
	args := s.Called()
	return args.Error(0)
}

func (s *SmtpConnectorMock) Disconnect() error {
	args := s.Called()
	return args.Error(0)
}

func (s *SmtpConnectorMock) Send(m *EmailMessage) error {
	args := s.Called()
	return args.Error(0)
}

func BenchmarkSendEmails(b *testing.B) {
	concurrencyLevels := []int{5, 50} //same time for benchmarks if lock mutex but faster if Read lock mutex in send (Send should not write then it's ok)
	for _, nconcurrent := range concurrencyLevels {
		b.Run(fmt.Sprintf("%d_clients", nconcurrent), func(b *testing.B) {
			connectorMock := &SimpleSmtpConnectorImpl{ErrorConnect: false, ErrorDisconnect: false, IsConnected: false, NSent: 0}
			emailsender := NewEmailSenderWithConnector(10, connectorMock)
			ch := make(chan struct{}, nconcurrent) //channel for limiting the number of concurrent jobs to nconcurrent
			wg := sync.WaitGroup{}
			for n := 0; n < 10; n++ {
				wg.Add(1)
				go func() {
					ch <- struct{}{}
					err := emailsender.Send(&EmailMessage{})
					//if an error occured (for example connection closed but trying to send) then fail the test
					if err != nil {
						logrus.Error(err)
					}
					<-ch
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}

/*
	In the use case of the application, a kafka consumer is a client and then creates an emailsender for each one.
	We are actually running n emailsenders that are sending b.N messages
*/
func BenchmarkSendEmailsUseCaseKafka(b *testing.B) {
	concurrencyLevels := []int{5, 50}
	for _, nconcurrent := range concurrencyLevels {
		b.Run(fmt.Sprintf("%d_clients", nconcurrent), func(b *testing.B) {
			wg := sync.WaitGroup{}
			for n := 0; n < nconcurrent; n++ {
				wg.Add(1)
				go func() {
					connectorMock := &SimpleSmtpConnectorImpl{ErrorConnect: false, ErrorDisconnect: false, IsConnected: false, NSent: 0}
					emailsender := NewEmailSenderWithConnector(10, connectorMock)
					for n := 0; n < b.N; n++ {
						err := emailsender.Send(&EmailMessage{})
						//if an error occured (for example connection closed but trying to send) then fail the test
						if err != nil {
							logrus.Error(err)
							b.Fail()
						}
					}
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}
