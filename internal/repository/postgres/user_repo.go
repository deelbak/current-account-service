package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"current-account-service/internal/models"
)

type userRepo struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) *userRepo { return &userRepo{db} }

func (r *userRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, role, login, password, created_at
		   FROM users WHERE login = $1`, login).
		Scan(&u.ID, &u.Name, &u.Role, &u.Login, &u.Password, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, role, login, password, created_at
		   FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Name, &u.Role, &u.Login, &u.Password, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

func (r *userRepo) GetAllByRole(ctx context.Context, role string) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, role, login, created_at FROM users WHERE role = $1`, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Role, &u.Login, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
