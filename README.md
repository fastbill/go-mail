# go-mail [![Build Status](https://travis-ci.org/fastbill/go-mail.svg?branch=master)](https://travis-ci.org/fastbill/go-mail) [![GoDoc](https://godoc.org/github.com/fastbill/go-mail?status.svg)](https://godoc.org/github.com/fastbill/go-mail)

Package mail provides an easy interface with interchangable backends.
Pull requests for additional mail providers are very welcome.

## Implementation roadmap

- [x] [Mandrill](https://mandrillapp.com)
- [ ] [net/smtp](https://golang.org/pkg/net/smtp/)

## Example usage

```go
package main

import (
	"github.com/fastbill/go-mail/v2"
	"github.com/fastbill/go-mail/v2/mandrill"
)

func main() {
	// Create and ping Mandrill client
	client := mandrill.MustNew("https://mandrillapp.com/api/1.0/", "my-token")
	err := client.Ping()
	if err != nil {
		panic(err)
	}

	// Send email
	err = client.Send(&mail.Mail{
		From:    &mail.Sender{Name: "FastBill GmbH", Email: "no-reply@fastbill.com"},
		To:      []mail.Recipient{mail.Recipient{Name: "Info", Email: "info@fastbill.com"}},
		Subject: "Hello world",
		HTML:    "<h1>Hello</h1>",
		Text:    "Hello",
		Options: &mandrill.Options{
			Important: true,
		},
	})
	if err != nil {
		panic(err)
	}
}
```
