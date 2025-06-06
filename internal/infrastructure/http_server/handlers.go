package http_server

import (
	"net/http"

	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
)

// Purchase handlers
func PurchaseDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement purchase download
	cookies.NotifySuccess(w, "Download iniciado")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func PurchaseCreateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement purchase creation
	cookies.NotifySuccess(w, "Compra iniciada")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Stripe handlers
func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Stripe checkout
	cookies.NotifySuccess(w, "Checkout iniciado")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Stripe webhook
	w.WriteHeader(http.StatusOK)
}

// Ebook handlers
func EbookIndexView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement ebook list view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func EbookCreateView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement ebook create view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func EbookCreateSubmit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement ebook create
	cookies.NotifySuccess(w, "Ebook criado com sucesso")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func EbookUpdateView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement ebook update view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func EbookShowView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement ebook show view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func EbookUpdateSubmit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement ebook update
	cookies.NotifySuccess(w, "Ebook atualizado com sucesso")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Send handler
func SendViewHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement send view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
