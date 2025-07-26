package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Initialize Stripe with API key
	stripe.Key = config.AppConfig.StripeSecretKey

	// Get session token from cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		log.Printf("Session token not found in cookie: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Não autorizado",
		})
		return
	}

	log.Printf("Session token found: %s", sessionCookie.Value)

	// Find user by session token
	userRepository := repository.NewGormUserRepository(database.DB)
	user := userRepository.FindBySessionToken(sessionCookie.Value)
	if user == nil {
		log.Printf("User not found for session token: %s", sessionCookie.Value)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Não autorizado",
		})
		return
	}

	log.Printf("User found: %s", user.Email)

	// Validate CSRF token
	csrfToken := r.Header.Get("X-CSRF-Token")
	log.Printf("CSRF token from header: %s", csrfToken)
	log.Printf("User CSRF token: %s", user.CSRFToken)

	if csrfToken == "" {
		log.Printf("Token CSRF não encontrado no header")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Token CSRF não encontrado",
		})
		return
	}

	if csrfToken != user.CSRFToken {
		log.Printf("Token CSRF inválido para o usuário %s. Token recebido: %s, Token esperado: %s",
			user.Email, csrfToken, user.CSRFToken)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Token CSRF inválido",
		})
		return
	}

	// Get user's subscription
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, gov.NewHubDevService())
	subscription, err := subscriptionService.FindByUserID(user.ID)
	if err != nil {
		log.Printf("Erro ao buscar subscription do usuário: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar pagamento",
		})
		return
	}

	if subscription == nil || subscription.StripeCustomerID == "" {
		log.Printf("Usuário %s não possui subscription ou ID do Stripe", user.Email)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar pagamento: Cliente não encontrado",
		})
		return
	}

	// Log Stripe configuration
	log.Printf("Stripe Secret Key: %s", config.AppConfig.StripeSecretKey)
	log.Printf("Stripe Price ID: %s", config.AppConfig.StripePriceID)

	if config.AppConfig.StripeSecretKey == "" {
		log.Printf("Stripe Secret Key não configurada")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro de configuração do Stripe",
		})
		return
	}

	if config.AppConfig.StripePriceID == "" {
		log.Printf("Stripe Price ID não configurado")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro de configuração do Stripe",
		})
		return
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(subscription.StripeCustomerID),
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

	log.Printf("Criando sessão do Stripe com os parâmetros: %+v", params)

	s, err := session.New(params)
	if err != nil {
		log.Printf("Erro ao criar sessão do Stripe: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar resposta: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao processar resposta",
		})
		return
	}
}

func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Configure webhook options to ignore API version mismatch
	opts := webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	}

	event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"), config.AppConfig.StripeWebhookSecret, opts)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing checkout session: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find subscription by Stripe customer ID
		subscriptionRepository := gorm.NewSubscriptionGormRepository()
		subscriptionService := service.NewSubscriptionService(subscriptionRepository, gov.NewHubDevService())
		subscription, err := subscriptionService.FindByStripeCustomerID(session.Customer.ID)
		if err != nil {
			log.Printf("Error finding subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if subscription == nil {
			log.Printf("Subscription not found for Stripe customer ID: %s", session.Customer.ID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Update subscription status
		err = subscriptionService.ActivateSubscription(subscription, session.Customer.ID, session.Subscription.ID)
		if err != nil {
			log.Printf("Error updating subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.updated":
		var stripeSubscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &stripeSubscription)
		if err != nil {
			log.Printf("Error parsing subscription: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find subscription by Stripe customer ID
		subscriptionRepository := gorm.NewSubscriptionGormRepository()
		subscriptionService := service.NewSubscriptionService(subscriptionRepository, gov.NewHubDevService())
		subscription, err := subscriptionService.FindByStripeCustomerID(stripeSubscription.Customer.ID)
		if err != nil {
			log.Printf("Error finding subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if subscription == nil {
			log.Printf("Subscription not found for Stripe customer ID: %s", stripeSubscription.Customer.ID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Update subscription status
		var endDate *time.Time
		if stripeSubscription.CurrentPeriodEnd > 0 {
			endDate = &time.Time{}
			*endDate = time.Unix(stripeSubscription.CurrentPeriodEnd, 0)
		}
		err = subscriptionService.UpdateSubscriptionStatus(subscription, string(stripeSubscription.Status), endDate)
		if err != nil {
			log.Printf("Error updating subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.deleted":
		var stripeSubscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &stripeSubscription)
		if err != nil {
			log.Printf("Error parsing subscription: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find subscription by Stripe customer ID
		subscriptionRepository := gorm.NewSubscriptionGormRepository()
		subscriptionService := service.NewSubscriptionService(subscriptionRepository, gov.NewHubDevService())
		subscription, err := subscriptionService.FindByStripeCustomerID(stripeSubscription.Customer.ID)
		if err != nil {
			log.Printf("Error finding subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if subscription == nil {
			log.Printf("Subscription not found for Stripe customer ID: %s", stripeSubscription.Customer.ID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Update subscription status
		var endDate *time.Time
		if stripeSubscription.CurrentPeriodEnd > 0 {
			endDate = &time.Time{}
			*endDate = time.Unix(stripeSubscription.CurrentPeriodEnd, 0)
		}
		err = subscriptionService.UpdateSubscriptionStatus(subscription, "canceled", endDate)
		if err != nil {
			log.Printf("Error updating subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
