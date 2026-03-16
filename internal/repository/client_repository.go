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
