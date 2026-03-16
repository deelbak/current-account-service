package dto

type CreateClientRequest struct {
	IIN       string `json:"iin"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Phone     string `json:"phone"`
}
