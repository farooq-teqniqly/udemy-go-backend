package main

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

func regexpExample() {
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	fmt.Println(validEmail.MatchString("farooq@foo.com"))
	fmt.Println(validEmail.MatchString("farooqfoo.com"))
	fmt.Println(validEmail.MatchString("farooqfoo@.com"))
	fmt.Println(validEmail.MatchString("farooq.foo.com"))
	fmt.Println(validEmail.MatchString("@f.com"))
	fmt.Println(validEmail.MatchString("f@f.co"))
	fmt.Println(validEmail.MatchString("farooq@foo.c"))
}

type Email struct {
	Sender      string
	Recipient   string
	Subject     string
	Body        string
	Attachments []string
}

func templateExample() {
	emailTemplate := `
	
	To: {{.Recipient}}
	From: {{.Sender}}
	Subject: {{.Subject}}
	
	Message:
		{{.Body}}

	{{if .Attachments}}
	Attachments:
	{{range .Attachments}}
		- {{.}}
	{{end}}
	{{else}}
	Email has no attachments.
	{{end}}
	`

	template, err := template.New("email").Parse(emailTemplate)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	email := Email{
		Sender:      "farooq@foo.com",
		Recipient:   "alice@foo.com",
		Subject:     "Welcome",
		Body:        "Welcome to the company. Glad to have you here!",
		Attachments: []string{"attachment1", "attachment2"},
	}

	// err = template.Execute(os.Stdout, email)

	output := strings.Builder{}

	err = template.Execute(&output, email)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(output.String())

}

func main() {
	// regexpExample()
	templateExample()
}
