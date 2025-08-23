package presentation

import (
	"log/slog"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/payment/business"
	"github.com/anglesson/simple-web-server/internal/payment/data"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type PaymentHandler struct {
	paymentService *business.PaymentServiceImpl
}

func NewPaymentHandler() *PaymentHandler {
	paymentRepository := data.NewPaymentRepository(database.DB)
	gateway := data.NewStripeService(config.AppConfig.StripeSecretKey)
	sellerAccountService := business.NewPaymentService(gateway, paymentRepository)

	return &PaymentHandler{
		sellerAccountService,
	}
}

func (h *PaymentHandler) CreateAccountForSeller(w http.ResponseWriter, r *http.Request) {
	// TODO: Find Creator by user ID
	sellerAccount, err := h.paymentService.CreateAccount("Anglesson", 1) // TODO: Inject User ID and Name
	if err != nil {
		slog.Error("Erro ao criar conta do vendedor: %v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	template.View(w, r, "dashboard", map[string]interface{}{
		"Flash": map[string]string{
			"Message": "Conta criada com sucesso!",
		},
		"SellerAccount": sellerAccount,
	}, "admin")
	return
}
