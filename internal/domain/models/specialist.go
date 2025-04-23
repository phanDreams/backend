package models

import "time"

type Specialist struct {
	ID int64          `json:"id"`
	Name string       `json:"name"`
	FamilyName string `json:"family_name"`
	Phone string      `json:"phone"`
	Email string      `json:"email"`
	Password string   `json:"-"` // store hashed password; omit from JSON responses
	IsBanned bool    `json:"is_banned"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive bool    `json:"is_active"`
	IsVerified bool  `json:"is_verified"`
	CreatedAt time.Time  `json:"created_at"`
}