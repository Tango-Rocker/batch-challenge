package business

import (
	"crypto/tls"
	"log/slog"

	gomail "gopkg.in/mail.v2"
)

type Mail struct {
	To      string
	Subject string
	Body    string
}

type Emailer interface {
	Send(Mail) error
}

type EmailService struct {
	from   string
	client *gomail.Dialer
	l      *slog.Logger
}

func NewEmailService(cfg *MailConfig, l *slog.Logger) *EmailService {
	d := gomail.NewDialer(
		cfg.Host,     //"smtp.gmail.com",
		cfg.Port,     //587,
		cfg.Account,  //"batchappdemo555@gmail.com",
		cfg.Password, //"obtp phtd dyse egat",
	)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &EmailService{
		from:   cfg.Account,
		client: d,
		l:      l.With(slog.String("service", "mail")),
	}
}

func (eSrv *EmailService) Send(msg Mail) {
	m := gomail.NewMessage()

	m.SetHeader("From", eSrv.from)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)

	err := eSrv.client.DialAndSend(m)
	if err != nil {
		eSrv.l.Error(err.Error())
	}
}
