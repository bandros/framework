package framework

import (
	"bytes"
	"errors"
	"gopkg.in/mailgun/mailgun-go.v1"
	"html/template"
)

type Email struct {
	template  []string
	subject   string
	recipient string
	sender    string
	body      string
	data      map[string]interface{}
}

func (e *Email) Template(template ...string) {
	e.template = template
}
func (e *Email) Set(subject, recipient string) {
	e.subject = subject
	e.recipient = recipient
}
func (e *Email) Sender(email string) {
	e.sender = email
}
func (e *Email) Data(data map[string]interface{}) {
	e.data = data
}
func authEmail(e *Email) error {
	if e.sender == "" {
		e.sender = Config("mailSenderMailgun")
	}

	if Config("domainMailgun") == "" {
		return errors.New("Check domainMailgun in config.yml")
	} else if Config("publicKeyMailgun") == "" {
		return errors.New("Check publicKeyMailgun in config.yml")
	} else if Config("keyMailgun") == "" {
		return errors.New("Check keyMailgun in config.yml")
	} else if Config("mailSenderMailgun") == "" {
		return errors.New("Check mailSenderMailgun in config.yml")
	} else if e.sender == "" {
		return errors.New("Sender email not set")
	} else if e.recipient == "" {
		return errors.New("recipient email not set")
	} else if e.subject == "" {
		return errors.New("subject not set")
	}
	return nil
}

func sendEmail(e *Email) (*mailgun.Message, mailgun.Mailgun, error) {
	err := authEmail(e)
	if err != nil {
		return nil, nil, err
	}
	mg := mailgun.NewMailgun(Config("domainMailgun"), Config("keyMailgun"), Config("publicKeyMailgun"))
	message := mg.NewMessage(e.sender, e.subject, "", e.recipient)
	return message, mg, nil
}

func (e *Email) Send() (string, error) {
	message, mg, err := sendEmail(e)
	if err != nil {
		return "", err
	}
	templateEmail := []string{}
	for _, v := range e.template {
		templateEmail = append(templateEmail, "./pages/email/"+v+".html")
	}
	html, err := templateEmailFunc(e.data, templateEmail...)

	if err != nil {
		return "", err
	}

	message.SetHtml(html)
	_, id, err := mg.Send(message)
	return id, err
	//return html,nil
}

func (e *Email) SendRaw(html string) (string, error) {
	message, mg, err := sendEmail(e)
	if err != nil {
		return "", err
	}
	message.SetHtml(html)
	_, id, err := mg.Send(message)
	return id, err
}

func templateEmailFunc(data map[string]interface{}, emailTemplate ...string) (content string, err error) {
	tmpl, err := template.ParseFiles(emailTemplate...)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
