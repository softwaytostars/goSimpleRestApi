package emails

import (
	"bytes"
	"fmt"
	"goapi/config"
	"goapi/emails"
	"goapi/kafka"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResourceEmails struct {
	emailKafkaProducer *kafka.EmailKafkaProducer
}

type formEmailBody struct {
	From        string                  `form:"from"`
	To          []string                `form:"to[]"`
	CC          []string                `form:"cc"`
	BCC         []string                `form:"bcc"`
	Subject     string                  `form:"subject"`
	TextBody    string                  `form:"textBody"`
	HtmlBody    string                  `form:"htmlBody"`
	Attachments []*multipart.FileHeader `form:"attachments[]""`
}

// Endpoint to Post messages to kafka
// @Summary  Post messages to kafka
// @Description  Post messages to kafka
// @Success 200 "OK"
// @Failure 500 {object} httputil.HTTPError
// @Router /emails [post]
func (r *ResourceEmails) sendEmail(c *gin.Context) {

	var form formEmailBody
	err := c.ShouldBind(&form)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Cannot deserialize formEmailBody [err=%s]", err)})
		return
	}

	emailMessage := emails.EmailMessage{
		From:        form.From,
		To:          form.To,
		CC:          form.CC,
		BCC:         form.BCC,
		Subject:     form.Subject,
		TextContent: form.TextBody,
		HtmlContent: form.HtmlBody,
		Attachments: make(map[string][]byte),
	}

	for _, attachment := range form.Attachments {

		src, err := attachment.Open()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Cannot open attachment [err=%s]", attachment.Filename)})
		}
		defer src.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, src); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Cannot copy attachment content [err=%s]", attachment.Filename)})
		}
		emailMessage.Attachments[attachment.Filename] = buf.Bytes()
	}

	err = r.emailKafkaProducer.ProduceEmails(emailMessage)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot post message [err=%s]", err)})
	} else {
		c.IndentedJSON(http.StatusOK, nil)
	}
}

// RegisterHandlers register all handlers for a router
func RegisterHandlers(r *gin.Engine, configuration *config.Config) {
	resource := ResourceEmails{emailKafkaProducer: kafka.NewEmailKafkaProducer(&configuration.KafkaConfig)}

	r.POST("/emails", resource.sendEmail)
}
