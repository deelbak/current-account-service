package service

import (
	"context"

	"current-account-service/internal/models"
	"current-account-service/internal/repository"
)

type TaskService struct {
	tasks repository.TaskRepository
}

func NewTaskService(tasks repository.TaskRepository) *TaskService {
	return &TaskService{tasks: tasks}
}

func (s *TaskService) ListMyTasks(ctx context.Context, userID int64) ([]models.Task, error) {
	return s.tasks.ListByAssignedTo(ctx, userID)
}
