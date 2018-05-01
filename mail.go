package mail

import (
	"io"
	"net/http"
)

// Client can be implemented for various services
type Client interface {
	// Send sends an email through the client
	Send(*Mail) error
	// Ping the address and validate key
	Ping() error
}

// Mail is a mail you can send
type Mail struct {
	// From is the name and email where this email is sent from
	From *Sender
	// To holds one or more recipients
	To []Recipient
	// Subject holds the mail subject
	Subject string
	// Headers are standard mail headers
	Headers map[string]interface{}
	// HTML will be displayed as html in the email
	HTML string
	// Text will be displayed as plain text in the email
	Text string
	// Options can be a custom interface provided by the client implementation
	Options interface{}
}

// Recipient is an email/name combination
type Recipient struct {
	Name  string
	Email string
}

// Sender holds data about the sender
type Sender Recipient

// HTTPClient is an interface for just the Post-method of *http.Client for easy mocking
type HTTPClient interface {
	// Post issues a POST to the specified URL
	Post(url string, contentType string, payload io.Reader) (*http.Response, error)
}
