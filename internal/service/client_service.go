package service

import (
	"current-account-service/internal/models"
	"current-account-service/internal/repository"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) GetClientByIIN(iin string) (*models.Client, error) {
	return s.repo.GetByIIN(iin)
}
