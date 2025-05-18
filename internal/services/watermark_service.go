package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func getFilename(original string) string {
	now := time.Now()

	// Formatar data e hora: YYYYMMDD_HHmmss
	timestamp := now.Format("20060102_150405")

	log.Println(timestamp)

	// Separar nome e extensão
	ext := filepath.Ext(original)
	name := strings.TrimSuffix(original, ext)

	// Juntar nome + timestamp + extensão
	newFileName := fmt.Sprintf("%s_%s%s", name, timestamp, ext)

	log.Println(newFileName)

	return newFileName
}

// ApplyWatermark aplica uma marca d'água ao PDF com as informações do usuário
func ApplyWatermark(inputPath, name, cpf, email string) (string, error) {
	outputPDF := getFilename(inputPath)

	// Tentar criar o diretório ./temp se ele não existir
	if err := os.MkdirAll("./temp", 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Configuração com opções de processamento
	conf := model.NewDefaultConfiguration()
	// conf.ValidationMode = model.ValidationRelaxed // Relaxar validação

	wm, errParse := pdfcpu.ParseTextWatermarkDetails(
		fmt.Sprintf("%s  - %s - %s", name, cpf, email),
		"font:Helvetica, points:20, pos:c, fillc:#000000, scale:1.0, rot:45, op:0.1",
		true,
		types.POINTS,
	)

	if errParse != nil {
		log.Fatal(errParse)
	}

	err := api.AddWatermarksFile(inputPath, outputPDF, nil, wm, conf)

	if err != nil {
		fmt.Println("Erro ao configurar o stamp:", err)
		return "", err
	}

	fmt.Println("Stamp aplicado com sucesso!")

	return outputPDF, nil
}
