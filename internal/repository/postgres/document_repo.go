package postgres

import (
	"context"
	"database/sql"

	"current-account-service/internal/models"
)

type documentRepo struct{ db *sql.DB }

func NewDocumentRepository(db *sql.DB) *documentRepo { return &documentRepo{db} }

func (r *documentRepo) Create(ctx context.Context, doc *models.Document) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO documents (application_id, file_path, doc_type)
		 VALUES ($1, $2, $3) RETURNING id, uploaded_at`,
		doc.ApplicationID, doc.FilePath, doc.DocType).
		Scan(&doc.ID, &doc.UploadedAt)
}

func (r *documentRepo) ListByApplicationID(ctx context.Context, appID int64) ([]models.Document, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, application_id, file_path, doc_type, uploaded_at
		   FROM documents WHERE application_id = $1 ORDER BY uploaded_at`, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var docs []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(&d.ID, &d.ApplicationID, &d.FilePath, &d.DocType, &d.UploadedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}
