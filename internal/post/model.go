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
