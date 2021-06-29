package monitor

import (
	"fmt"
	"net/smtp"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func SendMail(cfg *Config, subject string, msg string) error {
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	auth := smtp.PlainAuth("", cfg.Mail.From, cfg.Mail.Password, smtpServer.host)

	msgBody := fmt.Sprintf("Subject: %v\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n"+
		"%v\n", subject, msg)

	if err := smtp.SendMail(smtpServer.Address(), auth, cfg.Mail.From, cfg.Mail.To, []byte(msgBody)); err != nil {
		return err
	}

	return nil
}
