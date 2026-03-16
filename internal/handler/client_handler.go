package handler

import (
	"current-account-service/internal/service"
	"encoding/json"
	"net/http"
	"strings"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	iin := strings.TrimPrefix(r.URL.Path, "/clients/")

	client, err := h.service.GetClientByIIN(iin)
	if err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
