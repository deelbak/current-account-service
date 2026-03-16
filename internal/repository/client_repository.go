package repository

import (
	"current-account-service/internal/models"
	"database/sql"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) GetByIIN(iin string) (*models.Client, error) {
	query := `
	SELECT id, iin, first_name, last_name, birth_date, phone
	FROM clients
	WHERE iin = $1
	`

	row := r.db.QueryRow(query, iin)

	var client models.Client

	err := row.Scan(
		&client.ID,
		&client.IIN,
		&client.FirstName,
		&client.LastName,
		&client.BirthDate,
		&client.Phone,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) Create(client *models.Client) error {
	query := `
	INSERT INTO clients (iin, first_name, last_name, birth_date, phone)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		client.IIN,
		client.FirstName,
		client.LastName,
		client.BirthDate,
		client.Phone,
	).Scan(&client.ID, &client.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
