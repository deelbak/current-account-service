package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"current-account-service/internal/models"
	"current-account-service/internal/repository"
	"current-account-service/internal/sse"
	"current-account-service/pkg/statemachine"
)

type ApplicationService struct {
	repo        repository.ApplicationRepository
	docRepo     repository.DocumentRepository
	accountRepo repository.AccountRepository
	taskRepo    repository.TaskRepository
	sm          *statemachine.Machine
	hub         *sse.Hub
	notify      *NotificationService
}

func NewApplicationService(
	repo repository.ApplicationRepository,
	docRepo repository.DocumentRepository,
	accountRepo repository.AccountRepository,
	taskRepo repository.TaskRepository,
	sm *statemachine.Machine,
	hub *sse.Hub,
	notify *NotificationService,
) *ApplicationService {
	return &ApplicationService{
		repo:        repo,
		docRepo:     docRepo,
		accountRepo: accountRepo,
		taskRepo:    taskRepo,
		sm:          sm,
		hub:         hub,
		notify:      notify,
	}
}

func (s *ApplicationService) applyTransition(
	ctx context.Context,
	appID int64,
	event statemachine.Event,
	actorID int64,
	actorRole string,
	comment string,
) (*models.Application, error) {
	app, err := s.repo.GetByIDForUpdate(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("get application: %w", err)
	}

	nextState, err := s.sm.Transition(statemachine.State(app.Status), event)
	if err != nil {
		return nil, fmt.Errorf("invalid transition: %w", err)
	}

	prevState := app.Status
	app.Status = string(nextState)

	// Указатели для nullable полей
	var actorIDPtr *int64
	if actorID != 0 {
		actorIDPtr = &actorID
	}
	actorRolePtr := &actorRole
	commentPtr := &comment

	ev := models.ApplicationEvent{
		ApplicationID: appID,
		FromState:     prevState,
		ToState:       string(nextState),
		Event:         string(event),
		ActorID:       actorIDPtr,
		ActorRole:     actorRolePtr,
		Comment:       commentPtr,
		OccurredAt:    time.Now(),
	}

	if err := s.repo.UpdateStatusAndLogEvent(ctx, app, ev); err != nil {
		return nil, fmt.Errorf("save transition: %w", err)
	}

	// SSE broadcast
	go func() {
		s.hub.Publish <- sse.Message{
			ApplicationID: appID,
			State:         string(nextState),
			Event:         string(event),
			ActorID:       actorID,
		}
	}()

	// Уведомление
	go s.notify.Send(ev)

	return app, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, clientID, managerID int64) (*models.Application, error) {
	app := &models.Application{
		ClientID:  clientID,
		CreatedBy: managerID,
		Status:    string(statemachine.StateDraft),
	}
	if err := s.repo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}
	return s.applyTransition(ctx, app.ID, statemachine.EventCreate, managerID, "manager", "")
}

func (s *ApplicationService) GetWithEvents(ctx context.Context, id int64) (*models.ApplicationWithEvents, error) {
	return s.repo.GetWithEvents(ctx, id)
}

func (s *ApplicationService) Submit(ctx context.Context, appID, managerID int64) (*models.Application, error) {
	return s.applyTransition(ctx, appID, statemachine.EventSubmit, managerID, "manager", "")
}

func (s *ApplicationService) Approve(ctx context.Context, appID, headID int64, comment string) (*models.Application, error) {
	if _, err := s.applyTransition(ctx, appID, statemachine.EventApprove, headID, "ul_head", comment); err != nil {
		return nil, err
	}
	// Автооткрытие счёта
	app, err := s.applyTransition(ctx, appID, statemachine.EventOpenAcct, 0, "system", "auto")
	if err != nil {
		return nil, err
	}
	// Создаём запись счёта в БД
	_ = s.openAccount(ctx, app.ClientID)
	return app, nil
}

func (s *ApplicationService) RequestRevision(ctx context.Context, appID, headID int64, comment string) (*models.Application, error) {
	return s.applyTransition(ctx, appID, statemachine.EventRevision, headID, "ul_head", comment)
}

func (s *ApplicationService) Reject(ctx context.Context, appID, headID int64, comment string) (*models.Application, error) {
	return s.applyTransition(ctx, appID, statemachine.EventReject, headID, "ul_head", comment)
}

func (s *ApplicationService) UploadDocument(
	ctx context.Context,
	appID int64,
	docType string,
	filename string,
	file io.Reader,
) (*models.Document, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%d_%s", uploadDir, appID, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	doc := &models.Document{
		ApplicationID: appID,
		FilePath:      filePath,
		DocType:       docType,
	}
	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("save document: %w", err)
	}
	return doc, nil
}

func (s *ApplicationService) openAccount(ctx context.Context, clientID int64) error {
	number := fmt.Sprintf("KZ%020d", time.Now().UnixNano()%1_000_000_000_000_000_000)
	acc := &models.Account{
		ClientID:      clientID,
		AccountNumber: number,
	}
	return s.accountRepo.Create(ctx, acc)
}
