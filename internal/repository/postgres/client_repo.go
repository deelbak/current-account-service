package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"current-account-service/internal/models"
)

type clientRepo struct{ db *sql.DB }

func NewClientRepository(db *sql.DB) *clientRepo { return &clientRepo{db} }

func (r *clientRepo) GetByIIN(ctx context.Context, iin string) (*models.Client, error) {
	c := &models.Client{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, iin, first_name, last_name, birth_date, phone, created_at
		   FROM clients WHERE iin = $1`, iin).
		Scan(&c.ID, &c.IIN, &c.FirstName, &c.LastName, &c.BirthDate, &c.Phone, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("client not found")
	}
	return c, err
}

func (r *clientRepo) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	c := &models.Client{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, iin, first_name, last_name, birth_date, phone, created_at
		   FROM clients WHERE id = $1`, id).
		Scan(&c.ID, &c.IIN, &c.FirstName, &c.LastName, &c.BirthDate, &c.Phone, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("client not found")
	}
	return c, err
}
