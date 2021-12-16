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

type SmtpConnector interface {
	ConnectionIsOpen() bool
	Connect() error
	Disconnect() error
	Send(m *EmailMessage) error
}

type DefaultSmtpConnectorImpl struct {
	dialer       *gomail.Dialer
	senderCloser gomail.SendCloser
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func NewDefaultSmtpConnectorImpl(config *config.EmailServerConfig) *DefaultSmtpConnectorImpl {
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
	return &DefaultSmtpConnectorImpl{dialer: dialer}
}

func (c *DefaultSmtpConnectorImpl) ConnectionIsOpen() bool {
	return c.senderCloser != nil
}

func (c *DefaultSmtpConnectorImpl) Connect() error {

	senderCloser, err := c.dialer.Dial()
	if err != nil {
		logrus.Errorf("Cannot connect to stmp server %s", err)
	}
	c.senderCloser = senderCloser
	return err
}

func (c *DefaultSmtpConnectorImpl) Disconnect() error {
	if c.senderCloser == nil {
		return nil
	}
	logrus.Info("Closing connection to stmp server")
	if err := c.senderCloser.Close(); err != nil {
		logrus.Errorf("Cannot close connection to stmp server %s", err)
		return err
	}
	c.senderCloser = nil
	return nil
}

func (c *DefaultSmtpConnectorImpl) Send(m *EmailMessage) error {
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
	return gomail.Send(c.senderCloser, email)
}

type EmailSender struct {
	mutex         sync.RWMutex
	channel       chan struct{}
	connector     SmtpConnector
	timeoutIdleMs int
}

func NewEmailSenderWithConnector(timeoutIdleMs int, connector SmtpConnector) *EmailSender {
	emailSender := EmailSender{
		channel:       make(chan struct{}, 1),
		connector:     connector,
		timeoutIdleMs: timeoutIdleMs}

	//launch the daemon
	emailSender.run_daemon()

	return &emailSender
}

func NewEmailSender(config *config.EmailServerConfig) *EmailSender {
	return NewEmailSenderWithConnector(config.TimeoutIdleConnectionMs, NewDefaultSmtpConnectorImpl(config))
}

func (s *EmailSender) Close() {
	close(s.channel)    //close the deamon
	s.closeConnection() // close the connection
}

func (s *EmailSender) closeConnection() error {
	//prevent to close while someone is sending
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.connector.Disconnect()
}

func (s *EmailSender) openConnection() error {
	//prevent to send while the connection is creating
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.connector.Connect()
}

func (s *EmailSender) ConnectionIsOpen() bool {
	//prevent to send while the connection is creating
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.connector.ConnectionIsOpen()
}

//the daemon serves to close the smtp connection if no message were sent in the time range of 30 seconds
func (s *EmailSender) run_daemon() {
	go func() {

		if !s.ConnectionIsOpen() {
			s.openConnection()
		}

		for {

			select {
			case _, ok := <-s.channel:
				if !ok {
					return
				}
				if !s.ConnectionIsOpen() {
					s.openConnection()
				}

			// Close the connection to the SMTP server if no email was sent in
			// the last timeoutIdleMs.
			case <-time.After(time.Duration(s.timeoutIdleMs) * time.Millisecond):
				if s.ConnectionIsOpen() {
					s.closeConnection()
				}
			}
		}

	}()
}

func (s *EmailSender) Send(m *EmailMessage) error {

	go func() {
		//tell the daemon we are sending a message
		s.channel <- struct{}{}
	}()

	//the lock in openConnection must be taken from another go routine or oustide the lock in this go routine
	if !s.ConnectionIsOpen() {
		s.openConnection()
	}

	//daemon cannot close connection while the email is sent
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.connector.Send(m)
}
