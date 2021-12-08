package emails

import (
	"context"
	"encoding/json"
	"goapi/config"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type EmailKafkaConsumer struct {
	mutex       sync.RWMutex //use mutex so that to not interrupt any consuming process in progress if close command is called
	kafkaReader *kafka.Reader
	emailSender *EmailSender
}

func (r *EmailKafkaConsumer) CloseConsumer() {
	r.mutex.Lock()

	if err := r.kafkaReader.Close(); err != nil {
		logrus.Error("failed to close writer:", err)
	}
	r.emailSender.Close()

	r.mutex.Unlock()
}

func NewEmailKafkaConsumer(configKafka *config.KafkaServerConfig, configEmailServer *config.EmailServerConfig) *EmailKafkaConsumer {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{configKafka.Uri},
		GroupID:  "consumer-group-emails",
		Topic:    emailTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	emailConsumer := EmailKafkaConsumer{kafkaReader: kafkaReader, emailSender: NewEmailSender(configEmailServer)}
	go emailConsumer.consumeEmails() //start consuming as soon as created
	return &emailConsumer
}

func (r *EmailKafkaConsumer) readMessage() (*kafka.Message, error) {
	m, err := r.kafkaReader.ReadMessage(context.Background())
	for err != nil {
		logrus.Errorf("Break consumer, cannot read message %e, Will retry", err)
		m, err = r.kafkaReader.ReadMessage(context.Background())
	}
	return &m, err
}
func (r *EmailKafkaConsumer) consumeEmails() {
	for {
		m, _ := r.readMessage()
		r.mutex.RLock()
		defer r.mutex.RUnlock()

		//logrus.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		var email EmailMessage
		if err := json.Unmarshal(m.Value, &email); err != nil {
			logrus.Error("Cannot Unmarshal email %s", err)
		}

		logrus.Info("[EmailKafkaConsumer] Sending email")
		err := r.emailSender.Send(&email)
		if err != nil {
			logrus.Error("[EmailKafkaConsumer] Cannot send email %s", err)
		}
	}
}
