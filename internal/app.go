package internal

import (
	"github.com/anglesson/simple-web-server/internal/client"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"github.com/anglesson/simple-web-server/pkg/template"
)

// App representa a aplicação principal com todos os módulos
type App struct {
	ClientModule *client.Module
	// Outros módulos serão adicionados aqui
}

// Module representa um módulo completo com todas as suas dependências
type Module struct {
	Handler    interface{}
	Service    interface{}
	Repository interface{}
}

// NewApp cria uma nova instância da aplicação com todos os módulos
func NewApp(templateRenderer template.TemplateRenderer, flashMessageFactory web.FlashMessageFactory) *App {
	// Inicializar repositórios
	clientRepository := gorm.NewClientGormRepository()
	creatorRepository := gorm.NewCreatorRepository(database.DB)

	// Inicializar serviços externos
	commonRFService := gov.NewHubDevService()

	// Inicializar módulo client
	clientService := service.NewClientService(clientRepository, creatorRepository, commonRFService)

	// Criar creatorService temporariamente para o módulo client
	// TODO: Refatorar para injetar o creatorService como parâmetro
	creatorService := service.NewCreatorService(creatorRepository, commonRFService, nil, nil, nil)

	clientHandler := client.NewClientHandler(clientService, creatorService, flashMessageFactory, templateRenderer)

	clientModule := &client.Module{
		Handler:    clientHandler,
		Service:    clientService,
		Repository: clientRepository,
	}

	return &App{
		ClientModule: clientModule,
	}
}

// GetClientHandler retorna o handler do módulo client
func (app *App) GetClientHandler() *client.ClientHandler {
	return app.ClientModule.Handler
}

// GetClientService retorna o service do módulo client
func (app *App) GetClientService() service.ClientService {
	return app.ClientModule.Service.(service.ClientService)
}
