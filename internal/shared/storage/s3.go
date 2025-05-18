package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getConfig() aws.Config {
	cfg, err := awsCfg.LoadDefaultConfig(context.TODO(),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.AppConfig.S3AccessKey,
			config.AppConfig.S3SecretKey,
			"", // session token (deixe vazio se não estiver usando)
		)),
		awsCfg.WithRegion(config.AppConfig.S3Region), // ou sua região
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
		Bucket: aws.String(config.AppConfig.S3BucketName), // Substitua pelo nome do seu bucket
		Key:    aws.String(filename),                      // Caminho dentro do bucket
		Body:   file,
	})
	if err != nil {
		log.Fatalf("erro ao fazer upload: %v", err)
	}

	fmt.Println("Upload realizado com sucesso!")
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.AppConfig.S3BucketName, config.AppConfig.S3Region, filename)
	return fileURL, nil
}

func GenerateDownloadLink(filename string) string {
	cfg := getConfig()

	s3Client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(s3Client)

	// Cria o comando GetObject
	params := &s3.GetObjectInput{
		Bucket: aws.String(config.AppConfig.S3BucketName),
		Key:    aws.String(filename),
		// Opcional: Content-Disposition para forçar download com nome customizado
		// ResponseContentDisposition: aws.String("attachment; filename=\"arquivo.txt\""),
	}

	// Gera o link pré-assinado com tempo de expiração
	presignedURL, err := presigner.PresignGetObject(context.TODO(), params, func(opts *s3.PresignOptions) {})
	if err != nil {
		log.Fatalf("Erro ao gerar URL pré-assinada: %v", err)
	}

	return presignedURL.URL
}

func GetFile(filename string) (string, error) {
	cfg := getConfig()

	s3Client := s3.NewFromConfig(cfg)

	params := &s3.GetObjectInput{
		Bucket: aws.String(config.AppConfig.S3BucketName),
		Key:    aws.String(filename),
	}
	output, err := s3Client.GetObject(context.TODO(), params)
	if err != nil {
		log.Fatalf("Erro ao gerar URL pré-assinada: %v", err)
		return "", err
	}

	// Caminho local onde você quer salvar o arquivo
	localPath := "/tmp/" + filename

	// Criar arquivo local
	f, err := os.Create(localPath)
	if err != nil {
		log.Fatalf("erro ao criar arquivo local: %v", err)
	}
	defer f.Close()

	// Copiar conteúdo do S3 para o arquivo local
	_, err = io.Copy(f, output.Body)
	if err != nil {
		log.Fatalf("erro ao salvar conteúdo: %v", err)
	}

	return localPath, nil
}
