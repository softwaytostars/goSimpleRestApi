package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"goapi/config"
	"goapi/emails"

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

func (p *EmailKafkaProducer) ProduceEmails(email emails.EmailMessage) error {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(email)
	if err != nil {
		logrus.Error("Cannot encode struct to bytes")
		return err
	}

	err = p.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("KeyEmails"),
			Value: reqBodyBytes.Bytes(),
		})
	if err != nil {
		logrus.Errorf("[EmailKafkaProducer]%s", err)
		return err
	}
	logrus.Debug("message written")
	return nil
}
