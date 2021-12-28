package rabbitmq

import (
	"goapi/emails"
	"testing"
	"time"
)

/*
	These benchmarks are meant to work with kafka and rabbitmq servers
*/

func createRabbitConsumer() *EmailRabbitMQConsumer {
	connectorMock := &emails.SimpleSmtpConnectorImpl{
		ErrorConnect:    false,
		ErrorDisconnect: false,
		IsConnected:     false,
		NSent:           0}
	emailsender := emails.NewEmailSenderWithConnector(10, connectorMock)
	return NewRabbitMQConsumerWithEmailSender(emailsender)
}

type TestObserverEmail struct {
	Sent        chan int
	currentSent int
}

func (o *TestObserverEmail) OnEmailSent() {
	o.currentSent++
	o.Sent <- o.currentSent
	//logrus.Infof("Sent %d", o.currentSent)
}

func BenchmarkRabbitMqUse(b *testing.B) {

	producer := NewRabbitMQProducer()

	observer := TestObserverEmail{Sent: make(chan int, 1)}

	var rabbitConsumers []*EmailRabbitMQConsumer

	nconcurrent := 1
	for n := 0; n < nconcurrent; n++ {
		consumer := createRabbitConsumer()
		rabbitConsumers = append(rabbitConsumers, consumer)
		consumer.AddObserver(&observer)
		consumer.ConsumeEmails()
	}

	//logrus.Infof("b.N = %d", b.N)
	for n := 0; n < b.N; n++ {
		producer.publish(emails.EmailMessage{From: "from"})
	}

	loopWaitAllEvents(b, &observer)

	go func() {
		for _, v := range rabbitConsumers {
			v.CloseConsumer()
		}
	}()
	close(observer.Sent)
}

func loopWaitAllEvents(b *testing.B, observer *TestObserverEmail) {
	for {
		select {
		case current, ok := <-observer.Sent:
			if !ok {
				break
			}
			if current == b.N {
				return
			}
		case <-time.After(60 * time.Second):
			return
		}
	}
}
