package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"current-account-service/internal/models"
)

type applicationRepo struct{ db *sql.DB }

func NewApplicationRepository(db *sql.DB) *applicationRepo { return &applicationRepo{db} }

func (r *applicationRepo) Create(ctx context.Context, app *models.Application) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO applications (client_id, created_by, status)
		 VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		app.ClientID, app.CreatedBy, app.Status).
		Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

func (r *applicationRepo) GetByID(ctx context.Context, id int64) (*models.Application, error) {
	app := &models.Application{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, client_id, created_by, status, created_at, updated_at
		   FROM applications WHERE id = $1`, id).
		Scan(&app.ID, &app.ClientID, &app.CreatedBy, &app.Status, &app.CreatedAt, &app.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application not found")
	}
	return app, err
}

// GetByIDForUpdate — SELECT FOR UPDATE, вызывать внутри транзакции
func (r *applicationRepo) GetByIDForUpdate(ctx context.Context, id int64) (*models.Application, error) {
	app := &models.Application{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, client_id, created_by, status, created_at, updated_at
		   FROM applications WHERE id = $1 FOR UPDATE`, id).
		Scan(&app.ID, &app.ClientID, &app.CreatedBy, &app.Status, &app.CreatedAt, &app.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application not found")
	}
	return app, err
}

// UpdateStatusAndLogEvent — атомарно: обновить статус + записать событие в audit_log
func (r *applicationRepo) UpdateStatusAndLogEvent(
	ctx context.Context,
	app *models.Application,
	ev models.ApplicationEvent,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	app.UpdatedAt = time.Now()

	_, err = tx.ExecContext(ctx,
		`UPDATE applications SET status = $1, updated_at = $2 WHERE id = $3`,
		app.Status, app.UpdatedAt, app.ID)
	if err != nil {
		return fmt.Errorf("update application: %w", err)
	}

	// Пишем в application_events (state machine audit)
	err = tx.QueryRowContext(ctx,
		`INSERT INTO application_events
		   (application_id, from_state, to_state, event, actor_id, actor_role, comment, occurred_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		ev.ApplicationID, ev.FromState, ev.ToState, ev.Event,
		ev.ActorID, ev.ActorRole, ev.Comment, ev.OccurredAt).
		Scan(&ev.ID)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	// Также пишем в общий audit_log (уже существующий в схеме)
	action := fmt.Sprintf("%s → %s", ev.FromState, ev.ToState)
	details := ev.Comment
	_, err = tx.ExecContext(ctx,
		`INSERT INTO audit_log (application_id, user_id, action, details)
		 VALUES ($1, $2, $3, $4)`,
		ev.ApplicationID, ev.ActorID, action, details)
	if err != nil {
		return fmt.Errorf("insert audit_log: %w", err)
	}

	return tx.Commit()
}

func (r *applicationRepo) GetWithEvents(ctx context.Context, id int64) (*models.ApplicationWithEvents, error) {
	app, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, application_id, from_state, to_state, event,
		        actor_id, actor_role, comment, occurred_at
		   FROM application_events
		  WHERE application_id = $1
		  ORDER BY occurred_at ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &models.ApplicationWithEvents{Application: *app}
	for rows.Next() {
		var ev models.ApplicationEvent
		if err := rows.Scan(
			&ev.ID, &ev.ApplicationID, &ev.FromState, &ev.ToState, &ev.Event,
			&ev.ActorID, &ev.ActorRole, &ev.Comment, &ev.OccurredAt,
		); err != nil {
			return nil, err
		}
		result.Events = append(result.Events, ev)
	}
	return result, rows.Err()
}

func (r *applicationRepo) ListByStatus(ctx context.Context, status string) ([]models.Application, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, client_id, created_by, status, created_at, updated_at
		   FROM applications WHERE status = $1 ORDER BY created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []models.Application
	for rows.Next() {
		var a models.Application
		if err := rows.Scan(&a.ID, &a.ClientID, &a.CreatedBy, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}
