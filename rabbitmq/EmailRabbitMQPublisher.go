package rabbitmq

import (
	"bytes"
	"encoding/json"
	"goapi/emails"

	"github.com/sirupsen/logrus"
)

type EmailRabbitMQProducer struct {
	queueName string
}

func NewRabbitMQProducer() *EmailRabbitMQProducer {
	return &EmailRabbitMQProducer{GetQueueName(QUEUE_EMAIL)}
}

func (e *EmailRabbitMQProducer) publish(email emails.EmailMessage) error {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(email)
	if err != nil {
		logrus.Error("Cannot encode struct to bytes")
		return err
	}

	err = getRabbitConnectionProducer().publishToQueue(e.queueName, reqBodyBytes.Bytes())
	if err != nil {
		logrus.Errorf("[EmailRabbitMQProducer]%s", err)
		return err
	}
	logrus.Debug("message written")
	return nil
}
