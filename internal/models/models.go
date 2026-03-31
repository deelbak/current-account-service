package models

import "time"

// ─── Users ───────────────────────────────────────────────────────────────────

type User struct {
	ID        int64     `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Role      string    `db:"role"       json:"role"` // manager | ul_head
	Login     string    `db:"login"      json:"login"`
	Password  string    `db:"password"   json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ─── Clients ──────────────────────────────────────────────────────────────────

type Client struct {
	ID        int64      `db:"id"         json:"id"`
	IIN       string     `db:"iin"        json:"iin"`
	FirstName string     `db:"first_name" json:"first_name"`
	LastName  string     `db:"last_name"  json:"last_name"`
	BirthDate *time.Time `db:"birth_date" json:"birth_date,omitempty"`
	Phone     *string    `db:"phone"      json:"phone,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

// ─── Applications ─────────────────────────────────────────────────────────────

type Application struct {
	ID        int64     `db:"id"         json:"id"`
	ClientID  int64     `db:"client_id"  json:"client_id"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	Status    string    `db:"status"     json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ─── Documents ────────────────────────────────────────────────────────────────

type Document struct {
	ID            int64     `db:"id"             json:"id"`
	ApplicationID int64     `db:"application_id" json:"application_id"`
	FilePath      string    `db:"file_path"      json:"file_path"`
	DocType       string    `db:"doc_type"       json:"doc_type"`
	UploadedAt    time.Time `db:"uploaded_at"    json:"uploaded_at"`
}

// ─── Tasks ────────────────────────────────────────────────────────────────────

type Task struct {
	ID            int64     `db:"id"             json:"id"`
	ApplicationID int64     `db:"application_id" json:"application_id"`
	AssignedTo    int64     `db:"assigned_to"    json:"assigned_to"`
	TaskType      string    `db:"task_type"      json:"task_type"` // REVIEW | REVISION
	Status        string    `db:"status"         json:"status"`    // OPEN | DONE
	CreatedAt     time.Time `db:"created_at"     json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"     json:"updated_at"`
}

// ─── Accounts ─────────────────────────────────────────────────────────────────

type Account struct {
	ID            int64     `db:"id"             json:"id"`
	ClientID      int64     `db:"client_id"      json:"client_id"`
	AccountNumber string    `db:"account_number" json:"account_number"`
	OpenedAt      time.Time `db:"opened_at"      json:"opened_at"`
}

// ─── Audit Log ────────────────────────────────────────────────────────────────

type AuditLog struct {
	ID            int64     `db:"id"             json:"id"`
	ApplicationID *int64    `db:"application_id" json:"application_id,omitempty"`
	UserID        *int64    `db:"user_id"        json:"user_id,omitempty"`
	Action        string    `db:"action"         json:"action"`
	Timestamp     time.Time `db:"timestamp"      json:"timestamp"`
	Details       *string   `db:"details"        json:"details,omitempty"`
}

// ─── Application Events (state machine audit) ─────────────────────────────────

type ApplicationEvent struct {
	ID            int64     `db:"id"             json:"id"`
	ApplicationID int64     `db:"application_id" json:"application_id"`
	FromState     string    `db:"from_state"     json:"from_state"`
	ToState       string    `db:"to_state"       json:"to_state"`
	Event         string    `db:"event"          json:"event"`
	ActorID       *int64    `db:"actor_id"       json:"actor_id,omitempty"`
	ActorRole     *string   `db:"actor_role"     json:"actor_role,omitempty"`
	Comment       *string   `db:"comment"        json:"comment,omitempty"`
	OccurredAt    time.Time `db:"occurred_at"    json:"occurred_at"`
}

// ─── Request/Response DTOs ────────────────────────────────────────────────────

type CreateApplicationRequest struct {
	ClientID int64 `json:"client_id" validate:"required"`
}

type SubmitRevisionRequest struct {
	Comment string `json:"comment"`
}

type ApplicationWithEvents struct {
	Application
	Events []ApplicationEvent `json:"events"`
}
