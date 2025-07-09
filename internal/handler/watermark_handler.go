package handler

import (
	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/service"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func WatermarkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-app-key") != config.AppConfig.AppKey {
		http.Error(w, "Invalid app key", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Invalid content", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Erro ao obter arquivo: %v", err)
		http.Error(w, "Erro ao obter arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Criar arquivo temporário para salvar o upload
	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		log.Printf("Erro ao criar arquivo temporário: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	// Copiar o conteúdo do upload para o arquivo temporário
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Printf("Erro ao copiar arquivo: %v", err)
		http.Error(w, "Erro ao processar arquivo", http.StatusInternalServerError)
		return
	}
	tempFile.Close()

	// Aplicar marca d'água
	outputPath, err := service.ApplyWatermark(tempFile.Name(), content)
	if err != nil {
		log.Printf("Erro ao aplicar marca d'água: %v", err)
		http.Error(w, "Erro ao processar arquivo", http.StatusInternalServerError)
		return
	}
	defer os.Remove(outputPath)

	// Configurar cabeçalhos para download
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileHeader.Filename))
	w.Header().Set("Content-Type", "application/pdf")

	// Servir o arquivo
	http.ServeFile(w, r, outputPath)
}
