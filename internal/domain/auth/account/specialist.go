package account

import (
	"time"
)

type Specialist struct {
	ID int64          `json:"id"`
	Name string       `json:"name"`
	FamilyName string `json:"family_name"`
	Phone string      `json:"phone"`
	Email string      `json:"email"`
	PasswordHash string   `json:"-"` // store hashed password; omit from JSON responses
	IsBanned bool    `json:"is_banned"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive bool    `json:"is_active"`
	IsVerified bool  `json:"is_verified"`
	CreatedAt time.Time  `json:"created_at"`
}

func (s *Specialist) TableName() string   { return "specialists" }
func (s *Specialist) Columns() []string   { return []string{"name","family_name","phone","email","password_hash"} }
func (s *Specialist) Values() []interface{} {
    return []interface{}{&s.Name, &s.FamilyName, &s.Phone, &s.Email, &s.PasswordHash}
}
func (s *Specialist) SetID(id int64)        { s.ID = int64(id)}
func (s *Specialist) GetID() int64 {
	return s.ID
}

func (s *Specialist) SetPasswordHash(h string) { s.PasswordHash = h }
func (s *Specialist) GetPasswordHash() string {
	return s.PasswordHash
}

func (s *Specialist) GetEmail() string {
	return s.Email
}



