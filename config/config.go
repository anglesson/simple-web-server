package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfiguration struct {
	AppName            string
	AppMode            string
	Port               string
	MailHost           string
	MailPort           string
	MailUsername       string
	MailPassword       string
	MailAuth           string
	MailFromAddress    string
	MailFromName       string
	MailContactAddress string
	S3AccessKey        string
	S3SecretKey        string
	S3Region           string
}

var AppConfig AppConfiguration

func LoadConfigs() {
	if AppConfig.AppMode == "development" || AppConfig.AppMode == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Erro ao carregar o arquivo .env")
		}
	}
	AppConfig.AppName = GetEnv("APPLICATION_NAME", "Web App")
	AppConfig.Port = GetEnv("PORT", "8080")
	AppConfig.MailHost = GetEnv("MAIL_HOST", "sandbox.smtp.mailtrap.io")
	AppConfig.MailPort = GetEnv("MAIL_PORT", "2525")
	AppConfig.MailUsername = GetEnv("MAIL_USERNAME", "cc54bb91ec44b9")
	AppConfig.MailPassword = GetEnv("MAIL_PASSWORD", "fd9493e107213b")
	AppConfig.MailAuth = GetEnv("MAIL_AUTH", "PLAIN")
	AppConfig.MailFromAddress = GetEnv("MAIL_FROM_ADDRESS", "no-reply@simpleweb.com")
	AppConfig.MailContactAddress = GetEnv("MAIL_FROM_ADDRESS", "no-reply@simpleweb.com")
	AppConfig.S3AccessKey = GetEnv("S3_ACCESS_KEY", "")
	AppConfig.S3SecretKey = GetEnv("S3_SECRET_KEY", "")
	AppConfig.S3Region = GetEnv("S3_REGION", "sa-east-1")
}

func GetEnv(key, fallback string) string {
	env, exists := os.LookupEnv(key)
	if exists {
		return env
	}

	return fallback
}
