/*
Package mail provides an easy interface with interchangable backends.

	// Create and ping Mandrill mailer.
	mandrillMailer := mandrill.MustNew("https://mandrillapp.com/api/1.0/", "my-token")
	err := mandrillMailer.Ping()
	if err != nil {
		panic(err)
	}

	// Create template mailer.
	templateMailer := MustNewStandardTemplateMailer(mandrillMailer, "/templates/*.tmpl")

	// Configure email for sending.
	template := &Template{
		Data: map[string]interface{}{
			"Foo": 1234,
		},
		TextPath: "hello.text.tmpl",
		HTMLPath: "hello.html.tmpl",
	}
	config := &Config{
		From:    &Address{Name: "FastBill GmbH", Email: "no-reply@fastbill.com"},
		To:      []Address{Address{Name: "Info", Email: "info@fastbill.com"}},
		Subject: "Hello world",
		Options: &mandrill.Options{
			Important: true,
		},
	}

	// Send email.
	if err := templateMailer.Send(template, config); err != nil {
		panic(err)
	}
*/
package mail
