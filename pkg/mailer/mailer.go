package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"net/smtp"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	Password string
	sender   string
}

func New(sender string, password string) Mailer {

	return Mailer{
		sender:   sender,
		Password: password,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}, subject string) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	body := htmlBody.String()

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", m.sender, m.Password, "smtp.gmail.com"),
		m.sender, []string{recipient}, []byte(body))

	return err

}
