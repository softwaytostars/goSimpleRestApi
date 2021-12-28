package rabbitmq

import (
	"encoding/json"
	"goapi/emails"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type EmailRabbitMQConsumer struct {
	queueName     string
	rabbitChannel *amqp.Channel
	emailSender   *emails.EmailSender
	Observers     []emails.IObserverEmailSent
}

func NewRabbitMQConsumerWithEmailSender(emailSender *emails.EmailSender) *EmailRabbitMQConsumer {

	rabbitConnection := getRabbitConnectionConsumer()
	if rabbitConnection == nil {
		return nil
	}

	ch, err := rabbitConnection.createChannel()

	if err != nil {
		return nil
	}

	return &EmailRabbitMQConsumer{
		queueName:     GetQueueName(QUEUE_EMAIL),
		rabbitChannel: ch,
		emailSender:   emailSender,
	}
}

func (e *EmailRabbitMQConsumer) CloseConsumer() {
	e.rabbitChannel.Close()
}

func (e *EmailRabbitMQConsumer) ConsumeEmails() {
	msgs, err := e.rabbitChannel.Consume(
		e.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	if err != nil {
		logrus.Errorf("Cannot consume %s", err)
		return
	}

	go func() {
		err := e.readMessages(msgs)
		if err != nil {
			return
		}
	}()
}

func (r *EmailRabbitMQConsumer) AddObserver(observer emails.IObserverEmailSent) {
	if observer == nil {
		logrus.Info("observer is null")
	}
	if r.Observers == nil {
		logrus.Info(r.Observers)
	}
	r.Observers = append(r.Observers, observer)
}

func (r *EmailRabbitMQConsumer) notifyObservers() {
	for _, v := range r.Observers {
		v.OnEmailSent()
	}
}

func (r *EmailRabbitMQConsumer) readMessages(msgs <-chan amqp.Delivery) error {
	for msg := range msgs {
		var email emails.EmailMessage
		if err := json.Unmarshal(msg.Body, &email); err != nil {
			logrus.New().Errorf("Cannot Unmarshal email %s", err)
			return err
		}

		//logrus.Info("[EmailRabbitMQConsumer] Sending email")
		err := r.emailSender.Send(&email)
		if err != nil {
			logrus.Errorf("[EmailKafkaConsumer] Cannot send email %s", err)
			return err
		} else {
			r.notifyObservers()
		}
	}
	return nil
}
