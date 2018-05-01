/*
Package mail provides an easy interface with interchangable backends

Create and ping Mandrill client

	client := mandrill.MustNew("https://mandrillapp.com/api/1.0/", "my-token")
	err := client.Ping()
	if err != nil {
		panic(err)
	}

Send email

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

*/
package mail
