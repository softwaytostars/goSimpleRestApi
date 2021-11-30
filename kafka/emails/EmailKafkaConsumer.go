package emails

import (
	"context"
	"goapi/config"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type EmailKafkaConsumer struct {
	mutex       sync.RWMutex //use mutex so that to not interrupt any consuming process in progress if close command is called
	kafkaReader *kafka.Reader
}

func (r *EmailKafkaConsumer) CloseConsumer() {
	r.mutex.Lock()
	if err := r.kafkaReader.Close(); err != nil {
		logrus.Error("failed to close writer:", err)
	}
	r.mutex.Unlock()
}

func NewEmailKafkaConsumer(configuration *config.KafkaServerConfig) *EmailKafkaConsumer {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{configuration.Uri},
		GroupID:  "consumer-group-emails",
		Topic:    emailTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	emailConsumer := EmailKafkaConsumer{kafkaReader: kafkaReader}
	go emailConsumer.consumeEmails() //start consuming as soon as created
	return &emailConsumer
}

func (r *EmailKafkaConsumer) consumeEmails() {
	for {
		m, err := r.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			logrus.Errorf("Break consumer, cannot read message %e", err)
			break
		}
		r.mutex.Lock()
		logrus.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		r.mutex.Unlock()
	}
}
