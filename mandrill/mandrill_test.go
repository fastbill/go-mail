package mandrill

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/fastbill/go-mail"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("valid url", func(t *testing.T) {
		client, err := New("https://fastbill.com", "")
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("invalid url", func(t *testing.T) {
		client, err := New(":", "")
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestMustNew(t *testing.T) {
	t.Run("no panic on valid url", func(t *testing.T) {
		MustNew("https://fastbill.com", "")
	})

	t.Run("panics on valid url", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNew(":", "")
		})
	})
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Post(url string, contentType string, payload io.Reader) (*http.Response, error) {
	b, err := ioutil.ReadAll(payload)
	if err != nil {
		panic(err)
	}

	args := m.Called(url, contentType, string(b))
	res, _ := args.Get(0).(*http.Response)

	return res, args.Error(1)
}

func TestSend(t *testing.T) {
	client := MustNew("http://foo.bar", "foobar").(*Client)

	t.Run("happy path", func(t *testing.T) {

		// Set mock client
		mockClient := new(MockClient)
		client.httpClient = mockClient

		expectedBody := `{"key":"foobar","message":{"html":"","text":"World","subject":"Hello","from_email":"foo@domain.com","from_name":"Foo","to":[{"email":"bar@domain.com","name":"Bar","type":"to"}]}}
`
		mockClient.On("Post", "http://foo.bar/messages/send.json", "application/json", expectedBody).
			Return(&http.Response{StatusCode: 200}, nil)

		err := client.Send(&mail.Mail{
			From:    &mail.Sender{Name: "Foo", Email: "foo@domain.com"},
			To:      []mail.Recipient{mail.Recipient{Name: "Bar", Email: "bar@domain.com"}},
			Subject: "Hello",
			Text:    "World",
		})
		assert.NoError(t, err)
	})

	t.Run("http.Client error", func(t *testing.T) {
		// Set mock client
		mockClient := new(MockClient)
		client.httpClient = mockClient

		expectedBody := `{"key":"foobar","message":{"html":"","text":"World","subject":"Hello","from_email":"foo@domain.com","from_name":"Foo","to":[{"email":"bar@domain.com","name":"Bar","type":"to"}]}}
`

		mockClient.On("Post", "http://foo.bar/messages/send.json", "application/json", expectedBody).
			Return(nil, errors.New("Something is broken"))

		err := client.Send(&mail.Mail{
			From:    &mail.Sender{Name: "Foo", Email: "foo@domain.com"},
			To:      []mail.Recipient{mail.Recipient{Name: "Bar", Email: "bar@domain.com"}},
			Subject: "Hello",
			Text:    "World",
		})
		assert.Error(t, err)
	})

	t.Run("wrong status code", func(t *testing.T) {
		// Set mock client
		mockClient := new(MockClient)
		client.httpClient = mockClient

		expectedBody := `{"key":"foobar","message":{"html":"","text":"World","subject":"Hello","from_email":"foo@domain.com","from_name":"Foo","to":[{"email":"bar@domain.com","name":"Bar","type":"to"}]}}
`

		mockClient.On("Post", "http://foo.bar/messages/send.json", "application/json", expectedBody).
			Return(&http.Response{StatusCode: 400}, nil)

		err := client.Send(&mail.Mail{
			From:    &mail.Sender{Name: "Foo", Email: "foo@domain.com"},
			To:      []mail.Recipient{mail.Recipient{Name: "Bar", Email: "bar@domain.com"}},
			Subject: "Hello",
			Text:    "World",
		})
		assert.Error(t, err)
	})
}

func TestPing(t *testing.T) {
	client := MustNew("http://foo.bar", "foobar").(*Client)

	// Set mock client
	mockClient := new(MockClient)
	client.httpClient = mockClient

	expectedBody := `{"key":"foobar"}
`
	mockClient.On("Post", "http://foo.bar/users/ping.json", "application/json", expectedBody).
		Return(&http.Response{StatusCode: 200}, nil)

	err := client.Ping()
	assert.NoError(t, err)
}

func TestSendPayload(t *testing.T) {
	client := MustNew("http://foo.bar", "foobar").(*Client)
	client.httpClient = nil

	// nolint: bodyclose
	_, err := client.sendPayload("", &payload{
		Message: &message{
			Headers: map[string]interface{}{
				"hola": func() {
					// This should not work
				},
			},
		},
	})
	require.Error(t, err)
}
