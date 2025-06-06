package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/go-chi/chi/v5"
)

type ClientAPIHandler struct {
	clientUseCase *application.ClientUseCase
}

func NewClientAPIHandler(clientUseCase *application.ClientUseCase) *ClientAPIHandler {
	return &ClientAPIHandler{
		clientUseCase: clientUseCase,
	}
}

func (h *ClientAPIHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input application.CreateClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.clientUseCase.CreateClient(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *ClientAPIHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	var input application.UpdateClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.ID = uint(id)
	output, err := h.clientUseCase.UpdateClient(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *ClientAPIHandler) ImportClients(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	input := application.ImportClientsInput{
		File: fileBytes,
	}

	output, err := h.clientUseCase.ImportClients(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *ClientAPIHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	search := r.URL.Query().Get("search")

	input := application.ListClientsInput{
		Page:     page,
		PageSize: pageSize,
		Term:     search,
	}

	output, err := h.clientUseCase.ListClients(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
