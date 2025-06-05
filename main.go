package main

import (
	"net/http"

	client_application "github.com/anglesson/simple-web-server/internal/client/application"
	client_http "github.com/anglesson/simple-web-server/internal/client/infrastructure/http_server"
	client_persistence "github.com/anglesson/simple-web-server/internal/client/infrastructure/persistence/gorm"
	common_infrastructure "github.com/anglesson/simple-web-server/internal/common/infrastructure"
	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

func main() {
	// --- Infrastructure Layer Initialization ---
	db := client_persistence.NewGormClientRepository(database.DB)
	rfService := common_infrastructure.NewHubDevService(config.AppConfig.HubDesenvolvedorApi, config.AppConfig.HubDesenvolvedorToken)

	// --- Application Layer Initialization ---
	clientUseCase := client_application.NewClientUseCase(db, rfService)

	// --- HTTP Server Initialization ---
	router := client_http.NewRouter(clientUseCase)
	http.ListenAndServe(":8080", router)
}
