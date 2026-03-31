package service

import (
	"context"

	"current-account-service/internal/models"
	"current-account-service/internal/repository"
)

type ClientService struct {
	clients repository.ClientRepository
}

func NewClientService(clients repository.ClientRepository) *ClientService {
	return &ClientService{clients: clients}
}

func (s *ClientService) GetByIIN(ctx context.Context, iin string) (*models.Client, error) {
	return s.clients.GetByIIN(ctx, iin)
}
