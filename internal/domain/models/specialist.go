package models

type Specialist struct {
	ID int64          `json:"id"`
	Name string       `json:"name"`
	FamilyName string `json:"family_name"`
	Phone string      `json:"phone"`
	Email string      `json:"email"`
	Password string   `json:"-"` // store hashed password; omit from JSON responses
	Is_banned bool    `json:"is_banned"`
	Is_deleted bool   `json:"is_deleted"`
	Is_active bool    `json:"is_active"`
	Is_verified bool  `json:"is_verified"`
	CreatedAt string  `json:"created_at"`
}