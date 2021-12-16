package emails

import (
	"fmt"
	"goapi/config"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
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
	connectorMock := &SimpleSmtpConnectorImpl{errorConnect: false, errorDisconnect: false, isConnected: false, nSent: 0}
	emailsender := NewEmailSenderWithConnector(10, connectorMock)
	return NewEmailKafkaConsumerWithEmailSender(&config.KafkaServerConfig{Uri: getKafkaServerUri()}, emailsender)
}

type TestObserverEmail struct {
	Sent        chan int
	currentSent int
}

func (o *TestObserverEmail) OnEmailSent() {
	o.currentSent++
	o.Sent <- o.currentSent
	logrus.Infof("Sent %d", o.currentSent)
}

func Benchmark_KafkaUse(b *testing.B) {

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

	//kafka takes time to get ready, just make sure it is
	time.Sleep(10 * time.Millisecond)

	b.Run(fmt.Sprintf("%d_nMessages", 0), func(b *testing.B) {

		observer.currentSent = 0
		for n := 0; n < b.N; n++ {
			kafkaProducer.ProduceEmails(EmailMessage{From: "from"})
		}

		select {
		case current, ok := <-observer.Sent:
			if !ok {
				break
			}
			if current == b.N {
				break
			}
		case <-time.After(60 * time.Second):
			break
		}
	})
	go func() {
		for _, v := range kafkaConsumers {
			v.CloseConsumer()
		}
		kafkaProducer.Close()
	}()
	close(observer.Sent)
}
