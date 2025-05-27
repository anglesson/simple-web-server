package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := middlewares.Auth(r)
	if user == nil {
		log.Printf("Usuário não autenticado ao tentar criar sessão de checkout")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Não autorizado",
		})
		return
	}

	// Validate CSRF token
	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" {
		log.Printf("Token CSRF não encontrado no header")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Token CSRF não encontrado",
		})
		return
	}

	if csrfToken != user.CSRFToken {
		log.Printf("Token CSRF inválido para o usuário %s. Token recebido: %s, Token esperado: %s",
			user.Email, csrfToken, user.CSRFToken)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Token CSRF inválido",
		})
		return
	}

	if user.StripeCustomerID == "" {
		log.Printf("Usuário %s não possui ID do Stripe", user.Email)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar pagamento: Cliente não encontrado",
		})
		return
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(user.StripeCustomerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(config.AppConfig.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://" + r.Host + "/settings?success=true"),
		CancelURL:  stripe.String("http://" + r.Host + "/settings?canceled=true"),
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("Erro ao criar sessão do Stripe: %v", err)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar pagamento: " + err.Error(),
		})
		return
	}

	response := struct {
		URL string `json:"url"`
	}{
		URL: s.URL,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar resposta: %v", err)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar resposta",
		})
		return
	}
}
