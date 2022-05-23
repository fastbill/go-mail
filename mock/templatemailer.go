package mock

import (
	tmock "github.com/stretchr/testify/mock"

	"github.com/fastbill/go-mail/v3"
)

// TemplateMailer is a mock implementation of the mail.TemplateMailer interface.
type TemplateMailer struct {
	tmock.Mock
}

// Send is a mock implementation of mail.TemplateMailer#Send.
func (m *TemplateMailer) Send(template *mail.Template, config *mail.Config) error {
	args := m.Called(template, config)

	return args.Error(0)
}
