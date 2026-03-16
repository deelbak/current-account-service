package handler

import (
	"current-account-service/internal/dto"
	"current-account-service/internal/models"
	"current-account-service/internal/service"
	"encoding/json"
	"net/http"
	"time"

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

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req dto.CreateClientRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	birthDate, err := time.Parse("YYYY-MM-DD", req.BirthDate)
	if err != nil {
		http.Error(w, "Invalid birth date format", http.StatusBadRequest)
		return
	}

	client := models.Client{
		IIN:       req.IIN,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: birthDate,
		Phone:     req.Phone,
	}

	err = h.service.CreateClient(&client)
	if err != nil {
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
