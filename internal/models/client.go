package models

import "time"

type Client struct {
	ID        int
	IIN       string
	FirstName string
	LastName  string
	BirthDate time.Time
	Phone     string
	CreatedAt time.Time
}
