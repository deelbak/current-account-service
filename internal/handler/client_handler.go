package handler

import (
	"current-account-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	iin := chi.URLParam(r, "iin")

	client, err := h.service.GetClientByIIN(iin)
	if err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
