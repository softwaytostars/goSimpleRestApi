package kafka

import (
	"goapi/config"
	"goapi/emails"
	"os"
	"testing"
	"time"
)

func getKafkaServerUri() string {
	uri := os.Getenv("KAFKA_URI")
	if len(uri) <= 0 {
		uri = "localhost:9092"
	}
	return uri
}

/*
	These benchmarks are meant to work with kafka and rabbitmq servers
*/

func createKafkaConsumer() *EmailKafkaConsumer {
	connectorMock := &emails.SimpleSmtpConnectorImpl{
		ErrorConnect:    false,
		ErrorDisconnect: false,
		IsConnected:     false,
		NSent:           0}
	emailsender := emails.NewEmailSenderWithConnector(10, connectorMock)
	return NewEmailKafkaConsumerWithEmailSender(&config.KafkaServerConfig{Uri: getKafkaServerUri()}, emailsender)
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

func BenchmarkKafkaUse(b *testing.B) {

	kafkaProducer := NewEmailKafkaProducer(&config.KafkaServerConfig{Uri: getKafkaServerUri()})

	observer := TestObserverEmail{Sent: make(chan int, 1)}

	var kafkaConsumers []*EmailKafkaConsumer

	nconcurrent := 1
	for n := 0; n < nconcurrent; n++ {
		consumer := createKafkaConsumer()
		kafkaConsumers = append(kafkaConsumers, consumer)
		consumer.AddObserver(&observer)
		consumer.ConsumeEmails()
	}
	for n := 0; n < b.N; n++ {
		kafkaProducer.ProduceEmails(emails.EmailMessage{From: "from"})
	}

	loopWaitAllEvents(b, &observer)

	go func() {
		for _, v := range kafkaConsumers {
			v.CloseConsumer()
		}
		kafkaProducer.Close()
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
