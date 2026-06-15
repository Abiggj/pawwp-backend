package community

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AccountID    uuid.UUID  `gorm:"type:uuid;not null" json:"account_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	CoverageArea string     `json:"coverage_area"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type ChannelPost struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ChannelID uuid.UUID  `gorm:"type:uuid;not null" json:"channel_id"`
	AuthorID  uuid.UUID  `gorm:"type:uuid;not null" json:"author_id"`
	Type      string     `gorm:"not null" json:"type"`   // sos, adoption, donation
	Status    string     `gorm:"default:''" json:"status"` // e.g., urgent, resolved
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	MediaURL  string     `json:"media_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Subscription struct {
	AccountID uuid.UUID `gorm:"type:uuid;primaryKey" json:"account_id"`
	ChannelID uuid.UUID `gorm:"type:uuid;primaryKey" json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}
