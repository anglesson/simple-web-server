package service

import (
	"strconv"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/pkg/mail"
)

type EmailService struct {
	emailService *mail.EmailService
}

func NewEmailService() *EmailService {
	mailPort, _ := strconv.Atoi(config.AppConfig.MailPort)
	mailer := mail.NewGoMailer(
		config.AppConfig.MailHost,
		mailPort,
		config.AppConfig.MailUsername,
		config.AppConfig.MailPassword)

	emailService := mail.NewEmailService(mailer)

	return &EmailService{
		emailService: emailService,
	}
}

func (s *EmailService) SendPasswordResetEmail(name, email, resetToken string) error {
	resetLink := config.AppConfig.Host + ":" + config.AppConfig.Port + "/reset-password?token=" + resetToken

	s.emailService.SendPasswordResetEmail(name, email, resetLink)
	return nil
}
