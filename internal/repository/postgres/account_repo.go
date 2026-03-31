package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"current-account-service/internal/models"
)

type accountRepo struct{ db *sql.DB }

func NewAccountRepository(db *sql.DB) *accountRepo { return &accountRepo{db} }

func (r *accountRepo) Create(ctx context.Context, acc *models.Account) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO accounts (client_id, account_number)
		 VALUES ($1, $2) RETURNING id, opened_at`,
		acc.ClientID, acc.AccountNumber).
		Scan(&acc.ID, &acc.OpenedAt)
}

func (r *accountRepo) GetByClientID(ctx context.Context, clientID int64) (*models.Account, error) {
	acc := &models.Account{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, client_id, account_number, opened_at
		   FROM accounts WHERE client_id = $1`, clientID).
		Scan(&acc.ID, &acc.ClientID, &acc.AccountNumber, &acc.OpenedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	return acc, err
}
