package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"current-account-service/internal/models"
)

type taskRepo struct{ db *sql.DB }

func NewTaskRepository(db *sql.DB) *taskRepo { return &taskRepo{db} }

func (r *taskRepo) Create(ctx context.Context, task *models.Task) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (application_id, assigned_to, task_type, status)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		task.ApplicationID, task.AssignedTo, task.TaskType, task.Status).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *taskRepo) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	t := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, application_id, assigned_to, task_type, status, created_at, updated_at
		   FROM tasks WHERE id = $1`, id).
		Scan(&t.ID, &t.ApplicationID, &t.AssignedTo, &t.TaskType, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	return t, err
}

func (r *taskRepo) ListByAssignedTo(ctx context.Context, userID int64) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, application_id, assigned_to, task_type, status, created_at, updated_at
		   FROM tasks
		  WHERE assigned_to = $1 AND status = 'OPEN'
		  ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.ApplicationID, &t.AssignedTo, &t.TaskType, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *taskRepo) CloseOpenTasks(ctx context.Context, appID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET status = 'DONE', updated_at = now()
		  WHERE application_id = $1 AND status = 'OPEN'`, appID)
	return err
}

func (r *taskRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET status = $1, updated_at = now() WHERE id = $2`, status, id)
	return err
}
