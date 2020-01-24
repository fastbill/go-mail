package mandrill

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/fastbill/go-mail/v3"
)

var (
	sendPath = "/messages/send.json"
	pingPath = "/users/ping.json"
)

// Mandrill is a mandrill implementation of mailer interface.
type Mandrill struct {
	key        string
	baseURL    *url.URL
	httpClient mail.HTTPClient

	sendEndpoint string
	pingEndpoint string
}

// New creates a new Mandrill mailer with the given parameters.
func New(baseURL string, key string) (*Mandrill, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Mandrill{
		key:        key,
		baseURL:    parsedURL,
		httpClient: &http.Client{},

		sendEndpoint: parsedURL.String() + sendPath,
		pingEndpoint: parsedURL.String() + pingPath,
	}, nil
}

// MustNew returns a new Mandrill mailer or panics.
func MustNew(baseURL, key string) *Mandrill {
	mailer, err := New(baseURL, key)
	if err != nil {
		panic(err)
	}

	return mailer
}

// Send sends an email.
func (m *Mandrill) Send(mail *mail.Mail) error {
	msg := messageFromMail(mail)
	res, err := m.sendPayload(m.sendEndpoint, &payload{
		Key:     m.key,
		Message: msg,
	})
	if err != nil {
		return err
	}

	if res.Body != nil {
		return res.Body.Close()
	}

	return nil
}

// Ping returns an error if the key or endpoint are wrong.
func (m *Mandrill) Ping() error {
	res, err := m.sendPayload(m.pingEndpoint, &payload{
		Key: m.key,
	})
	if err != nil {
		return err
	}

	if res.Body != nil {
		return res.Body.Close()
	}

	return nil
}

func (m *Mandrill) sendPayload(endpoint string, data *payload) (*http.Response, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, err
	}

	res, err := m.httpClient.Post(endpoint, "application/json", buf)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code: %d", res.StatusCode)
	}

	return res, nil
}

type payload struct {
	Key     string   `json:"key"`
	Message *message `json:"message,omitempty"`
	// TODO: Add more payload fields
}

type message struct {
	*Options

	HTML      string                 `json:"html"`
	Text      string                 `json:"text"`
	Subject   string                 `json:"subject"`
	FromEmail string                 `json:"from_email"`
	FromName  string                 `json:"from_name"`
	To        []recipient            `json:"to"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
}

func messageFromMail(m *mail.Mail) *message {
	options, _ := m.Options.(*Options)

	var to []recipient
	for _, r := range m.To {
		to = append(to, recipient{
			Email: r.Email,
			Name:  r.Name,
			Type:  "to",
		})
	}

	return &message{
		Options: options,

		FromEmail: m.From.Email,
		FromName:  m.From.Name,
		To:        to,
		Subject:   m.Subject,
		HTML:      m.HTML,
		Text:      m.Text,
		Headers:   m.Headers,
	}
}

// Options hold mandrill-specific data
// NOTE: Disabled maligned, because that's how the docs sort it
// nolint: maligned
type Options struct {
	Important               bool                   `json:"important,omitempty"`
	TrackOpens              bool                   `json:"track_opens,omitempty"`
	TrackClicks             bool                   `json:"track_clicks,omitempty"`
	AutoText                bool                   `json:"auto_text,omitempty"`
	AutoHTML                bool                   `json:"auto_html,omitempty"`
	InlineCSS               bool                   `json:"inline_css,omitempty"`
	URLStripQS              bool                   `json:"url_strip_qs,omitempty"`
	PreserveRecipients      bool                   `json:"preserve_recipients,omitempty"`
	ViewContentLink         bool                   `json:"view_content_link,omitempty"`
	BCCAddress              string                 `json:"bcc_address,omitempty"`
	TrackingDomain          string                 `json:"tracking_domain,omitempty"`
	SigningDomain           string                 `json:"signing_domain,omitempty"`
	ReturnPathDomain        string                 `json:"return_path_domain,omitempty"`
	Merge                   bool                   `json:"merge,omitempty"`
	MergeLanguage           string                 `json:"merge_language,omitempty"`
	GlobalMergeVars         []MergeVar             `json:"global_merge_vars,omitempty"`
	MergeVars               []MergeVars            `json:"merge_vars,omitempty"`
	Tags                    []string               `json:"tags,omitempty"`
	SubAccount              string                 `json:"sub_account,omitempty"`
	GoogleAnalyticsDomains  []string               `json:"google_analytics_domains,omitempty"`
	GoogleAnalyticsCampaign string                 `json:"google_analytics_campaign,omitempty"`
	Metadata                map[string]interface{} `json:"metadata,omitempty"`
	RecipientMetadata       *RecipientMetadata     `json:"recipient_metadata,omitempty"`
	Attachments             []Attachment           `json:"attachments,omitempty"`
	Images                  []Attachment           `json:"images,omitempty"`
}

// RecipientMetadata holds metadata about the recipient
type RecipientMetadata struct {
	Rcpt   string                 `json:"rcpt,omitempty"`
	Values map[string]interface{} `json:"values,omitempty"`
}

// MergeVars holds merge vars
type MergeVars struct {
	Rcpt string     `json:"rcpt,omitempty"`
	Vars []MergeVar `json:"vars,omitempty"`
}

// MergeVar holds global merge vars
type MergeVar struct {
	Name    string      `json:"name,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

// Attachment can be any file to attach to the message
type Attachment struct {
	Mime    string    `json:"type,omitempty"`
	Name    string    `json:"name,omitempty"`
	Content io.Reader `json:"content,omitempty"`
}

type recipient struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}
