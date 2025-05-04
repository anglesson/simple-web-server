package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	app_config "github.com/anglesson/simple-web-server/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Upload(file io.Reader, filename string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			app_config.AppConfig.S3SecretKey,
			app_config.AppConfig.S3AccessKey,
			"", // session token (deixe vazio se não estiver usando)
		)),
		config.WithRegion(app_config.AppConfig.S3Region), // ou sua região
	)

	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	// Instancia o client S3
	s3Client := s3.NewFromConfig(cfg)

	// Envia o arquivo
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("docffy-homolog"), // Substitua pelo nome do seu bucket
		Key:    aws.String(filename),         // Caminho dentro do bucket
		Body:   file,
	})
	if err != nil {
		log.Fatalf("erro ao fazer upload: %v", err)
	}

	fmt.Println("Upload realizado com sucesso!")
	return nil
}
