package mail

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) Send(eml *Mail) error {
	args := m.Called(eml)
	return args.Error(0)
}

func (m *MockMailer) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewStandardTemplateMailer(t *testing.T) {
	t.Run("invalid template path", func(t *testing.T) {
		mockMailer := &MockMailer{}

		_, err := NewStandardTemplateMailer(mockMailer, "./invalid/*.tmpl")

		require.Error(t, err)
		assert.EqualError(t, err, "template: pattern matches no files: `./invalid/*.tmpl`")
	})

	t.Run("success", func(t *testing.T) {
		mockMailer := &MockMailer{}

		_, err := NewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		require.NoError(t, err)
	})
}

func TestMustNewStandardTemplateMailer(t *testing.T) {
	mockMailer := &MockMailer{}

	assert.Panics(t, func() {
		MustNewStandardTemplateMailer(mockMailer, "./invalid/*.tmpl")
	})
}

func TestStandardTemplateMailer_Send(t *testing.T) {
	t.Run("config is required", func(t *testing.T) {
		mockMailer := &MockMailer{}
		templateMailer := MustNewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		template := &Template{}
		err := templateMailer.Send(template, nil)

		require.Error(t, err)
		assert.EqualError(t, err, "config parameter is required")
	})

	t.Run("template is required", func(t *testing.T) {
		mockMailer := &MockMailer{}
		templateMailer := MustNewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		config := &Config{}
		err := templateMailer.Send(nil, config)

		require.Error(t, err)
		assert.EqualError(t, err, "template parameter is required")
	})

	t.Run("failed to execute plain text template", func(t *testing.T) {
		mockMailer := &MockMailer{}
		templateMailer := MustNewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		template := &Template{
			TextPath: "foo_txt.tmpl",
		}
		err := templateMailer.Send(template, &Config{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), `template: no template "foo_txt.tmpl"`)
	})

	t.Run("failed to execute HTML template", func(t *testing.T) {
		mockMailer := &MockMailer{}
		templateMailer := MustNewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		template := &Template{
			TextPath: "foo_text.tmpl",
			HTMLPath: "foo_htm.tmpl",
		}
		err := templateMailer.Send(template, &Config{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), `template: no template "foo_htm.tmpl"`)
	})

	t.Run("failed to send email", func(t *testing.T) {
		mockMailer := &MockMailer{}
		templateMailer := MustNewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		errUnexpected := errors.New("unexpected error")
		mockMailer.On("Send", mock.Anything).Return(errUnexpected)

		template := &Template{
			TextPath: "foo_text.tmpl",
			HTMLPath: "foo_html.tmpl",
		}
		err := templateMailer.Send(template, &Config{})

		mockMailer.AssertExpectations(t)
		assert.Equal(t, errUnexpected, err)
	})

	t.Run("success", func(t *testing.T) {
		mockMailer := &MockMailer{}
		templateMailer := MustNewStandardTemplateMailer(mockMailer, "./testdata/*.tmpl")

		config := &Config{
			From: &Address{
				Name:  "joe",
				Email: "joedoe@example.com",
			},
			To: []Address{
				{
					Name:  "jane",
					Email: "janedoe@example.com",
				},
			},
			Subject: "hello",
			Options: nil,
		}
		eml := &Mail{
			Config: *config,
			HTML:   "foohtml\n",
			Text:   "footext\n",
		}
		mockMailer.On("Send", eml).Return(nil)

		template := &Template{
			Data: map[string]interface{}{
				"FooVar": "foo",
			},
			TextPath: "foo_text.tmpl",
			HTMLPath: "foo_html.tmpl",
		}
		err := templateMailer.Send(template, config)

		mockMailer.AssertExpectations(t)
		require.NoError(t, err)
	})
}
