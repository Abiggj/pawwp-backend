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

type MedicalRecord struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PetID            uuid.UUID `gorm:"type:uuid;not null;index" json:"pet_id"`
	AccountID        uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	RecordType       string    `gorm:"size:50;not null" json:"record_type"` // vaccination, medication, visit, allergy, surgery, lab_result
	Title            string    `gorm:"size:255;not null" json:"title"`
	Veterinarian     string    `gorm:"size:255" json:"veterinarian"`
	DateAdministered string    `gorm:"size:50" json:"date_administered"`
	ExpiryDate       string    `gorm:"size:50" json:"expiry_date"`
	Dosage           string    `gorm:"size:100" json:"dosage"`
	Notes            string    `gorm:"type:text" json:"notes"`
	DocumentURL      string    `gorm:"size:500" json:"document_url"`
	Status           string    `gorm:"size:50;default:'active'" json:"status"` // up_to_date, due_soon, expired, active, completed
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
