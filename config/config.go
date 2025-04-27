package config

import "os"

type AppConfiguration struct {
	AppName            string
	Port               string
	MailHost           string
	MailPort           string
	MailUsername       string
	MailPassword       string
	MailAuth           string
	MailFromAddress    string
	MailFromName       string
	MailContactAddress string
}

var AppConfig AppConfiguration

func LoadConfigs() {
	AppConfig.AppName = GetEnv("APPLICATION_NAME", "Web App")
	AppConfig.Port = GetEnv("PORT", "8080")
	AppConfig.MailHost = GetEnv("MAIL_HOST", "sandbox.smtp.mailtrap.io")
	AppConfig.MailPort = GetEnv("MAIL_PORT", "2525")
	AppConfig.MailUsername = GetEnv("MAIL_USERNAME", "cc54bb91ec44b9")
	AppConfig.MailPassword = GetEnv("MAIL_PASSWORD", "fd9493e107213b")
	AppConfig.MailAuth = GetEnv("MAIL_AUTH", "PLAIN")
	AppConfig.MailFromAddress = GetEnv("MAIL_FROM_ADDRESS", "no-reply@simpleweb.com")
	AppConfig.MailFromAddress = GetEnv("MAIL_FROM_ADDRESS", "no-reply@simpleweb.com")
}

func GetEnv(key, fallback string) string {
	env, exists := os.LookupEnv(key)
	if exists {
		return env
	}

	return fallback
}
