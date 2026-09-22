package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type         string     `gorm:"size:50;not null;default:'parent'" json:"type"` // parent, shelter_admin, volunteer, shelter
	Name         string     `gorm:"size:255" json:"name"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string     `json:"-"`
	RefreshToken string     `json:"-"`
	ShelterName  string     `gorm:"size:255" json:"shelter_name"`
	RoleTitle    string     `gorm:"size:100" json:"role_title"`
	Bio          string     `json:"bio"`
	IsVerified   bool       `json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
