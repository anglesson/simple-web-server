package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	app_config "github.com/anglesson/simple-web-server/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			app_config.AppConfig.S3AccessKey,
			app_config.AppConfig.S3SecretKey,
			"", // session token (deixe vazio se não estiver usando)
		)),
		config.WithRegion(app_config.AppConfig.S3Region), // ou sua região
	)
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}
	return cfg
}

func Upload(file io.Reader, filename string) (string, error) {
	cfg := getConfig()

	// Instancia o client S3
	s3Client := s3.NewFromConfig(cfg)

	// Envia o arquivo
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(app_config.AppConfig.S3BucketName), // Substitua pelo nome do seu bucket
		Key:    aws.String(filename),                          // Caminho dentro do bucket
		Body:   file,
	})
	if err != nil {
		log.Fatalf("erro ao fazer upload: %v", err)
	}

	fmt.Println("Upload realizado com sucesso!")
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", app_config.AppConfig.S3BucketName, app_config.AppConfig.S3Region, filename)
	return fileURL, nil
}

func GenerateDownloadLink(filename string) string {
	cfg := getConfig()

	s3Client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(s3Client)

	// Cria o comando GetObject
	params := &s3.GetObjectInput{
		Bucket: aws.String(app_config.AppConfig.S3BucketName),
		Key:    aws.String(filename),
		// Opcional: Content-Disposition para forçar download com nome customizado
		// ResponseContentDisposition: aws.String("attachment; filename=\"arquivo.txt\""),
	}

	// Gera o link pré-assinado com tempo de expiração
	presignedURL, err := presigner.PresignGetObject(context.TODO(), params, func(opts *s3.PresignOptions) {
		opts.Expires = time.Minute * 15 // link expira em 15 minutos
	})
	if err != nil {
		log.Fatalf("Erro ao gerar URL pré-assinada: %v", err)
	}

	return presignedURL.URL
}
