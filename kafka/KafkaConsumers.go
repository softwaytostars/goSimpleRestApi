package kafka

import (
	"goapi/config"
	"goapi/kafka/emails"
	"sync"

	"github.com/sirupsen/logrus"
)

type TypeConsumer int

const (
	EmailConsumer TypeConsumer = iota
	Sms
)

type KafkaConsumers struct {
	config    *config.Config
	consumers map[TypeConsumer][]Consumer
}

var instance *KafkaConsumers
var once sync.Once

func GetInstanceKafkaConsumers(config *config.Config) *KafkaConsumers {
	once.Do(func() {
		instance = &KafkaConsumers{config: config, consumers: make(map[TypeConsumer][]Consumer)}
	})
	return instance
}

func (c *KafkaConsumers) StartConsumers(n int, consumerType TypeConsumer) {
	for i := 0; i < n; i++ {
		switch consumerType {
		case EmailConsumer:
			c.consumers[EmailConsumer] = append(c.consumers[EmailConsumer], emails.NewEmailKafkaConsumer(&c.config.KafkaConfig, &c.config.EmailServerConfig))
		default:
			logrus.Error("Not supported consumer type")
		}
	}
}

func (c *KafkaConsumers) retrieveListOfConsumers(consumerType TypeConsumer) []Consumer {
	switch consumerType {
	case EmailConsumer:
		return c.consumers[EmailConsumer]
	default:
		logrus.Error("Not supported consumer type")
		return nil
	}
}

func (c *KafkaConsumers) StopConsumers(n int, consumerType TypeConsumer) {
	liste := c.retrieveListOfConsumers(consumerType)
	if liste != nil {
		//remove n first elements
		for _, consumer := range liste[0:n] {
			consumer.CloseConsumer()
		}
	}
}
