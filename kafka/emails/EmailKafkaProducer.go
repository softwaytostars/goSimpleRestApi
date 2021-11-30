package emails

import (
	"context"
	"goapi/config"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type EmailKafkaProducer struct {
	kafkaWriter *kafka.Writer
}

func NewEmailKafkaProducer(configuration *config.KafkaServerConfig) *EmailKafkaProducer {
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(configuration.Uri),
		Topic:    emailTopic,
		Balancer: &kafka.LeastBytes{},
	}
	return &EmailKafkaProducer{kafkaWriter}
}

func (p *EmailKafkaProducer) Close() {
	if err := p.kafkaWriter.Close(); err != nil {
		logrus.Error("failed to close writer:", err)
	}
}

func (p *EmailKafkaProducer) ProduceEmails() {
	err := p.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key1"),
			Value: []byte("one!"),
		},
		kafka.Message{
			Key:   []byte("Key2"),
			Value: []byte("two!"),
		},
		kafka.Message{
			Key:   []byte("Key3"),
			Value: []byte("three!"),
		})
	if err != nil {
		logrus.Error(err)
	}
}
