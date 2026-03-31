package repository

import (
	"context"
	"current-account-service/internal/models"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetAllByRole(ctx context.Context, role string) ([]models.User, error)
}

type ClientRepository interface {
	GetByIIN(ctx context.Context, iin string) (*models.Client, error)
	GetByID(ctx context.Context, id int64) (*models.Client, error)
}

type ApplicationRepository interface {
	Create(ctx context.Context, app *models.Application) error
	GetByID(ctx context.Context, id int64) (*models.Application, error)
	GetByIDForUpdate(ctx context.Context, id int64) (*models.Application, error) // SELECT FOR UPDATE
	UpdateStatusAndLogEvent(ctx context.Context, app *models.Application, ev models.ApplicationEvent) error
	GetWithEvents(ctx context.Context, id int64) (*models.ApplicationWithEvents, error)
	ListByStatus(ctx context.Context, status string) ([]models.Application, error)
}

type DocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) error
	ListByApplicationID(ctx context.Context, appID int64) ([]models.Document, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	ListByAssignedTo(ctx context.Context, userID int64) ([]models.Task, error)
	CloseOpenTasks(ctx context.Context, appID int64) error
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type AccountRepository interface {
	Create(ctx context.Context, acc *models.Account) error
	GetByClientID(ctx context.Context, clientID int64) (*models.Account, error)
}

type AuditRepository interface {
	Log(ctx context.Context, entry *models.AuditLog) error
}
