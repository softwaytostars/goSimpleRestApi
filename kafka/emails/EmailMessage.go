package emails

import (
	"io/ioutil"
	"path/filepath"
)

type EmailMessage struct {
	From        string
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	TextContent string
	HtmlContent string
	Attachments map[string][]byte
}

func (m *EmailMessage) AddAttachment(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *EmailMessage) SendTo(to []string) {
	m.To = to
}
