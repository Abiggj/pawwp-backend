package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type         string     `json:"type"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Bio          string     `json:"bio"`
	IsVerified   bool       `json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
