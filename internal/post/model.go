package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PetID       uuid.UUID  `gorm:"type:uuid;not null" json:"pet_id"`
	MediaURL    string     `json:"media_url"`
	Caption     string     `json:"caption"`
	IsShowcased bool       `gorm:"default:false" json:"is_showcased"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Boop struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_account_post" json:"account_id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_account_post" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Woof struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type BoopDetail struct {
	ID         uuid.UUID `json:"id"`
	PostID     uuid.UUID `json:"post_id"`
	AccountID  uuid.UUID `json:"account_id"`
	UserName   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar"`
	UserType   string    `json:"user_type"`
	BoopedAt   time.Time `json:"booped_at"`
}

type WoofDetail struct {
	ID           uuid.UUID `json:"id"`
	PostID       uuid.UUID `json:"post_id"`
	AccountID    uuid.UUID `json:"account_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	AuthorType   string    `json:"author_type"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type PostWithInteractions struct {
	Post
	BoopsCount int64        `json:"boops_count"`
	HasBooped  bool         `json:"has_booped"`
	Boops      []BoopDetail `json:"boops"`
	WoofsCount int64        `json:"woofs_count"`
	Woofs      []WoofDetail `json:"woofs"`
	ViewsCount int64        `json:"views_count"`
}

type PostView struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PostID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"post_id"`
	AccountID    *uuid.UUID `gorm:"type:uuid;index" json:"account_id,omitempty"`
	ViewerKey    string     `gorm:"size:100;not null;index" json:"viewer_key"`
	DwellSeconds int        `gorm:"default:2" json:"dwell_seconds"`
	IsAuthentic  bool       `gorm:"default:true" json:"is_authentic"`
	CreatedAt    time.Time  `json:"created_at"`
}
