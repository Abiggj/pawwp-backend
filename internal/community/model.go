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

type ChannelAdmin struct {
	ChannelID uuid.UUID `gorm:"type:uuid;primaryKey" json:"channel_id"`
	AccountID uuid.UUID `gorm:"type:uuid;primaryKey" json:"account_id"`
	Role      string    `gorm:"default:'admin'" json:"role"` // owner, medical_admin, dispatch_admin, admin
	AdminName string    `gorm:"-" json:"admin_name,omitempty"`
	AdminEmail string   `gorm:"-" json:"admin_email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TownhallQuestion struct {
	ID           uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AuthorID     uuid.UUID        `gorm:"type:uuid;not null" json:"author_id"`
	AuthorName   string           `json:"author_name"`
	AuthorType   string           `json:"author_type"` // parent, shelter
	Title        string           `gorm:"not null" json:"title"`
	Content      string           `gorm:"not null" json:"content"`
	Category     string           `gorm:"default:'general'" json:"category"` // urgent_health, diet, behavior, general
	IsUrgent     bool             `gorm:"default:false" json:"is_urgent"`
	HelpfulVotes int              `gorm:"default:0" json:"helpful_votes"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Answers      []TownhallAnswer `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}

type TownhallAnswer struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuestionID   uuid.UUID `gorm:"type:uuid;not null" json:"question_id"`
	AuthorID     uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	AuthorName   string    `json:"author_name"`
	AuthorType   string    `json:"author_type"` // parent, shelter
	Content      string    `gorm:"not null" json:"content"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"` // shelter/vet response
	HelpfulVotes int       `gorm:"default:0" json:"helpful_votes"`
	CreatedAt    time.Time `json:"created_at"`
}
