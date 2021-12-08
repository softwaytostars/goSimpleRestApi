package emails

import (
	"crypto/tls"
	"goapi/config"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	mutex        sync.RWMutex
	channel      chan bool
	dialer       *gomail.Dialer
	senderCloser gomail.SendCloser
}

func NewEmailSender(config *config.EmailServerConfig) *EmailSender {
	ch := make(chan bool)
	var dialer *gomail.Dialer
	if config.UseStartTLS {
		dialer = gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	} else {
		dialer = &gomail.Dialer{
			Host:      config.Host,
			Port:      config.Port,
			Username:  config.Username,
			Password:  config.Password,
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
			SSL:       false,
		}
	}

	emailSender := EmailSender{channel: ch, dialer: dialer}
	emailSender.run_daemon()
	return &emailSender
}

func (s *EmailSender) Close() {
	close(s.channel)    //close the deamon
	s.closeConnection() // close the connection
}

func (s *EmailSender) closeConnection() {
	if s.senderCloser == nil {
		return
	}
	//prevent to close while someone is sending
	s.mutex.Lock()
	defer s.mutex.Unlock()

	logrus.Info("Closing connection to stmp server")
	if err := s.senderCloser.Close(); err != nil {
		logrus.Errorf("Cannot close connection to stmp server %s", err)
	}
	s.senderCloser = nil
}

func (s *EmailSender) openConnection() error {
	//prevent to send while the connection is creating
	s.mutex.Lock()
	defer s.mutex.Unlock()

	senderCloser, err := s.dialer.Dial()
	if err != nil {
		logrus.Errorf("Cannot connect to stmp server %s", err)
	}
	s.senderCloser = senderCloser
	return err
}

//the daemon serves to close the smtp connection if no message were sent in the time range of 30 seconds
func (s *EmailSender) run_daemon() {
	go func() {
		open := s.openConnection() != nil
		for {
			select {
			case _, ok := <-s.channel:
				if !ok {
					return
				}
				if !open {
					open = s.openConnection() != nil
				}

				// Close the connection to the SMTP server if no email was sent in
				// the last 30 seconds.
			case <-time.After(30 * time.Second):
				s.closeConnection()
				open = false
			}
		}

	}()
}

func (s *EmailSender) Send(m *EmailMessage) error {
	email := gomail.NewMessage()
	email.SetHeader("From", m.From)
	email.SetHeader("To", m.To...)
	email.SetHeader("Cc", m.CC...)
	email.SetHeader("Subject", m.Subject)
	if len(m.TextContent) > 0 && len(m.HtmlContent) > 0 {
		email.SetBody("text/html", m.HtmlContent)
		email.AddAlternative("text/plain", m.TextContent)
	} else if len(m.TextContent) > 0 {
		email.SetBody("text/plain", m.TextContent)
	} else {
		email.SetBody("text/html", m.HtmlContent)
	}

	for filename, content := range m.Attachments {
		email.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(content)
			return err
		}))
	}

	//tell the daemon we are sending a message
	s.channel <- true
	//open the connection if not already open
	if s.senderCloser == nil {
		err := s.openConnection()
		if err != nil {
			return err
		}
	}

	//deamon cannot close connection while the email is sent
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	//send the email
	return gomail.Send(s.senderCloser, email)
}
