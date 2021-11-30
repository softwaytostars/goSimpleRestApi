package emails

import (
	"goapi/config"
	"goapi/kafka/emails"

	"github.com/gin-gonic/gin"
)

type ResourceEmails struct {
	emailKafkaProducer *emails.EmailKafkaProducer
}

func (r *ResourceEmails) sendEmail(c *gin.Context) {
	r.emailKafkaProducer.ProduceEmails()
}

// RegisterHandlers register all handlers for a router
func RegisterHandlers(r *gin.Engine, configuration *config.Config) {
	resource := ResourceEmails{emailKafkaProducer: emails.NewEmailKafkaProducer(&configuration.KafkaConfig)}

	r.POST("/emails", resource.sendEmail)
}
