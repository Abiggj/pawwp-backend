package pet

import (
	"time"

	"github.com/google/uuid"
)

type Pet struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AccountID      uuid.UUID  `gorm:"type:uuid;not null" json:"account_id"`
	Username       string     `gorm:"uniqueIndex;size:30;not null" json:"username"`
	Name           string     `json:"name"`
	Species        string     `json:"species"`
	Breed          string     `json:"breed"`
	Bio            string     `json:"bio"`
	AdoptionStatus string     `json:"adoption_status"`
	IsVerified     bool       `json:"is_verified"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
