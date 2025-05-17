package mail

import (
	"fmt"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/models"
)

type EmailService struct {
	mailer Mailer
}

func NewEmailService(mailer Mailer) *EmailService {
	return &EmailService{
		mailer: mailer,
	}
}

func (s *EmailService) SendPasswordResetEmail(name, email string, resetLink string) {
	data := map[string]interface{}{
		"ResetLink": resetLink,
		"Name":      name,
		"Title":     "Recover your password!",
	}

	s.mailer.From(config.AppConfig.MailFromAddress)
	s.mailer.To(email)
	s.mailer.Subject("Recover your password!")
	s.mailer.Body(NewEmail("reset_password", data))
	s.mailer.Send()
}

func (s *EmailService) SendAccountConfirmation(name, email, token string) {
	data := map[string]interface{}{
		"Name":               name,
		"Title":              "Confirm your account!",
		"AppName":            config.AppConfig.AppName,
		"Contact":            config.AppConfig.MailFromAddress,
		"ConfirmAccountLink": "/account-confirmation?token=" + token + "&name=" + name + "&email=" + email,
	}

	s.mailer.From(config.AppConfig.MailFromAddress)
	s.mailer.To(email)
	s.mailer.Subject("Confirm your account")
	s.mailer.Body(NewEmail("account_confirmation", data))
	s.mailer.Send()
}

func (s *EmailService) SendLinkToDownload(purchases []*models.Purchase) {
	for _, purchase := range purchases {
		data := map[string]interface{}{
			"Name":              purchase.Client.Name,
			"Title":             "Seu ebook chegou!",
			"AppName":           config.AppConfig.AppName,
			"Contact":           config.AppConfig.MailFromAddress,
			"EbookDownloadLink": fmt.Sprintf("/download/ebook/%d", purchase.EbookID),
		}

		s.mailer.From(config.AppConfig.MailFromAddress)
		s.mailer.To(purchase.Client.Contact.Email)
		s.mailer.Subject("Seu ebook chegou!")
		s.mailer.Body(NewEmail("ebook_download", data))
		s.mailer.Send()
	}
}
