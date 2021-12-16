package emails

import (
	"context"
	"encoding/json"
	"goapi/config"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type EmailKafkaConsumer struct {
	kafkaReader *kafka.Reader
	emailSender *EmailSender
	Observers   []IObserverEmailSent
}

func (r *EmailKafkaConsumer) CloseConsumer() {
	if err := r.kafkaReader.Close(); err != nil {
		logrus.Error("failed to close writer:", err)
	}
	r.emailSender.Close()
}

func NewEmailKafkaConsumer(configKafka *config.KafkaServerConfig, configEmailServer *config.EmailServerConfig) *EmailKafkaConsumer {
	consumer := NewEmailKafkaConsumerWithEmailSender(configKafka, NewEmailSender(configEmailServer))
	consumer.ConsumeEmails()
	return consumer
}

func NewEmailKafkaConsumerWithEmailSender(configKafka *config.KafkaServerConfig, emailSender *EmailSender) *EmailKafkaConsumer {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{configKafka.Uri},
		GroupID:  "consumer-group-emails",
		Topic:    emailTopic,
		MinBytes: 0,    // 1B
		MaxBytes: 10e6, // 10MB
	})
	emailConsumer := EmailKafkaConsumer{
		kafkaReader: kafkaReader,
		emailSender: emailSender,
	}

	return &emailConsumer
}

func (r *EmailKafkaConsumer) AddObserver(observer IObserverEmailSent) {
	r.Observers = append(r.Observers, observer)
}

func (r *EmailKafkaConsumer) notifyObservers() {
	for _, v := range r.Observers {
		v.OnEmailSent()
	}
}

func (r *EmailKafkaConsumer) readMessage() (*kafka.Message, error) {
	m, err := r.kafkaReader.ReadMessage(context.Background())
	/*
		for err != nil {
			logrus.Errorf("Break consumer, cannot read message %e, Will retry", err)
			m, err = r.kafkaReader.ReadMessage(context.Background())
		}
	*/
	return &m, err
}
func (r *EmailKafkaConsumer) ConsumeEmails() {
	go func() {
		for {
			err := r.readMessages()
			if err != nil {
				//logrus.Errorf("Consumer is stopping to fetch messages %s", err)
				return
			}
		}
	}()
}

func (r *EmailKafkaConsumer) readMessages() error {
	m, err := r.readMessage()
	if err != nil {
		//logrus.Errorf("[EmailKafkaConsumer] readMessage %s", err)
		return err
	}
	//logrus.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	var email EmailMessage
	if err = json.Unmarshal(m.Value, &email); err != nil {
		logrus.New().Errorf("Cannot Unmarshal email %s", err)
		return err
	}

	logrus.Info("[EmailKafkaConsumer] Sending email")
	err = r.emailSender.Send(&email)
	if err != nil {
		logrus.Errorf("[EmailKafkaConsumer] Cannot send email %s", err)
	} else {
		r.notifyObservers()
	}
	return err
}
