package mail

import (
	"bytes"
	"io"
	"net/http"
	"text/template"
)

// Mailer can be implemented for various services.
type Mailer interface {
	// Send sends an email.
	Send(*Mail) error

	// Ping pings the address.
	Ping() error
}

// TemplateMailer is a representation of mailer which compiles templates to email body (HTML and plain text).
type TemplateMailer interface {
	Send(template *Template, config *Config) error
}

// HTTPClient is an interface for just the Post-method of *http.Mailer for easy mocking.
type HTTPClient interface {
	// Post issues a POST to the specified URL.
	Post(url string, contentType string, payload io.Reader) (*http.Response, error)
}

// Config is a configurable part of an email.
type Config struct {
	// From is the name and email where this email is sent from.
	From *Address

	// To holds one or more recipients.
	To []Address

	// Subject holds the email subject.
	Subject string

	// Headers are standard email headers.
	Headers map[string]interface{}

	// Options can be a custom interface provided by the mailer implementation.
	Options interface{}
}

// Mail is a mail you can send.
type Mail struct {
	Config

	// HTML will be displayed as html in the email.
	HTML string

	// Text will be displayed as plain text in the email.
	Text string
}

// Address is an email/name combination.
type Address struct {
	Name  string
	Email string
}

// Template holds all info needed for templates to compile to an email body.
type Template struct {
	Data     map[string]interface{}
	TextPath string
	HTMLPath string
}

// StandardTemplateMailer is a standard (default) implementation of template mailer interface.
type StandardTemplateMailer struct {
	mailer   Mailer
	template *template.Template
}

// NewStandardTemplateMailer creates a new StandardTemplateMailer mailer.
func NewStandardTemplateMailer(mailer Mailer, templatePath string) (*StandardTemplateMailer, error) {
	tpl, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, err
	}

	return &StandardTemplateMailer{
		mailer:   mailer,
		template: tpl,
	}, nil
}

// MustStandardTemplateMailer returns a new StandardTemplateMailer mailer or panics.
func MustStandardTemplateMailer(mailer Mailer, templatePath string) *StandardTemplateMailer {
	tpl, err := NewStandardTemplateMailer(mailer, templatePath)
	if err != nil {
		panic(err)
	}

	return tpl
}

// Send sends an email.
func (m *StandardTemplateMailer) Send(template *Template, config *Config) error {
	textBuf := &bytes.Buffer{}
	if err := m.template.ExecuteTemplate(textBuf, template.TextPath, template.Data); err != nil {
		return err
	}

	htmlBuf := &bytes.Buffer{}
	if err := m.template.ExecuteTemplate(htmlBuf, template.HTMLPath, template.Data); err != nil {
		return err
	}

	return m.mailer.Send(&Mail{
		Config: *config,
		HTML:   htmlBuf.String(),
		Text:   textBuf.String(),
	})
}
