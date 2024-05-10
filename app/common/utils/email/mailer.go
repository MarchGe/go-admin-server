package email

import (
	"fmt"
	"github.com/wneessen/go-mail"
)

type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	From     string `mapstructure:"from"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseSSL   bool   `mapstructure:"useSSL"`
}

type Mailer struct {
	cfg *MailConfig
}

func (m *Mailer) createClient() (*mail.Client, error) {
	client, err := mail.NewClient(m.cfg.Host, mail.WithPort(m.cfg.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(m.cfg.Username), mail.WithPassword(m.cfg.Password))
	if err != nil {
		return nil, err
	}
	client.SetSSL(m.cfg.UseSSL)
	return client, nil
}

func (m *Mailer) Send(msg ...*mail.Msg) error {
	client, err := m.createClient()
	if err != nil {
		return fmt.Errorf("mailer createClient error, %w", err)
	}
	err = client.DialAndSend(msg...)
	if err != nil {
		return fmt.Errorf("mailer send error, %w", err)
	}
	return nil
}
