package rabbitmq

import (
	"errors"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type EnumTypeQueue int

const (
	QUEUE_EMAIL EnumTypeQueue = iota
)

func GetQueueName(typeQueue EnumTypeQueue) string {
	switch typeQueue {
	case QUEUE_EMAIL:
		return "emails_queue"
	default:
		return ""
	}
}

type RabbitMQConnection struct {
	connection *amqp.Connection
}

func getRabbitMQServerUri() string {
	uri := os.Getenv("RABBITMQ_URI")
	if len(uri) <= 0 {
		uri = "localhost:5672"
	}
	return uri
}

var onceProducer sync.Once
var onceConsumer sync.Once
var instanceProducer *RabbitMQConnection
var instanceConsumer *RabbitMQConnection

func createRabbitConnection(doOnSuccess func(*amqp.Connection)) *RabbitMQConnection {
	instance := &RabbitMQConnection{}
	connection, err := amqp.Dial("amqp://guest:guest@" + getRabbitMQServerUri() + "/")
	if err != nil {
		logrus.Error(err)
	} else {
		instance.connection = connection
		doOnSuccess(connection)
	}
	return instance
}

func getRabbitConnectionConsumer() *RabbitMQConnection {
	onceConsumer.Do(func() {
		instanceConsumer = createRabbitConnection(createRabbitMQArchitecture)

	})
	return instanceConsumer
}

func getRabbitConnectionProducer() *RabbitMQConnection {
	onceProducer.Do(func() {
		instanceProducer = createRabbitConnection(createRabbitMQArchitecture)
	})
	return instanceProducer
}

func (r *RabbitMQConnection) createChannel() (*amqp.Channel, error) {
	if r.connection == nil {
		return nil, errors.New("no connection")
	}
	return r.connection.Channel()
}

func createRabbitMQArchitecture(connection *amqp.Connection) {
	ch, err := connection.Channel()
	if err != nil {
		logrus.Error(err)
	} else {
		ch.ExchangeDeclare(
			"live_exchange", // name
			"direct",        // type
			true,            // durable
			false,           // auto-deleted
			false,           // internal
			false,           // no-wait
			nil,             // arguments
		)
		ch.QueueDeclare(GetQueueName(QUEUE_EMAIL), true, false, false, false, nil)
	}
	ch.Close()
}

func (r *RabbitMQConnection) publishToQueue(queueName string, message []byte) error {
	if instanceProducer == nil || instanceProducer.connection == nil {
		logrus.Error("Cannot publish")
		return errors.New("no connection")
	}

	ch, err := instanceProducer.connection.Channel()
	if err != nil {
		logrus.Error("Cannot create channel")
		return errors.New("no channel")
	}
	defer ch.Close()

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})

	if err != nil {
		logrus.Errorf("Cannot publish %s", err)
	}
	return err
}

func (r *RabbitMQConnection) Close() error {
	return r.connection.Close()
}
