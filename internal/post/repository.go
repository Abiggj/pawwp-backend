package post

import (
	"time"

	"github.com/Abiggj/pawwp/internal/database"
	"github.com/google/uuid"
)

type Repository interface {
	Create(post *Post) error
	FindByID(id string) (*Post, error)
	Get48hFeed() ([]Post, error)
	GetShowcaseByPetID(petID uuid.UUID) ([]Post, error)
	CountShowcaseByPetID(petID uuid.UUID) (int64, error)
	GetArchiveByPetID(petID uuid.UUID) ([]Post, error)
	Update(post *Post) error

	AddBoop(boop *Boop) error
	RemoveBoop(accountID, postID uuid.UUID) error
	AddWoof(woof *Woof) error
	GetWoofsByPostID(postID uuid.UUID) ([]Woof, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(post *Post) error {
	return database.DB.Create(post).Error
}

func (r *repository) FindByID(id string) (*Post, error) {
	var p Post
	err := database.DB.Where("id = ?", id).First(&p).Error
	return &p, err
}

func (r *repository) Get48hFeed() ([]Post, error) {
	var posts []Post
	fortyEightHoursAgo := time.Now().Add(-48 * time.Hour)
	err := database.DB.Where("created_at >= ?", fortyEightHoursAgo).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *repository) GetShowcaseByPetID(petID uuid.UUID) ([]Post, error) {
	var posts []Post
	err := database.DB.Where("pet_id = ? AND is_showcased = ?", petID, true).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *repository) CountShowcaseByPetID(petID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&Post{}).Where("pet_id = ? AND is_showcased = ?", petID, true).Count(&count).Error
	return count, err
}

func (r *repository) GetArchiveByPetID(petID uuid.UUID) ([]Post, error) {
	var posts []Post
	err := database.DB.Where("pet_id = ?", petID).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *repository) Update(post *Post) error {
	return database.DB.Save(post).Error
}

func (r *repository) AddBoop(boop *Boop) error {
	return database.DB.Create(boop).Error
}

func (r *repository) RemoveBoop(accountID, postID uuid.UUID) error {
	return database.DB.Delete(&Boop{}, "account_id = ? AND post_id = ?", accountID, postID).Error
}

func (r *repository) AddWoof(woof *Woof) error {
	return database.DB.Create(woof).Error
}

func (r *repository) GetWoofsByPostID(postID uuid.UUID) ([]Woof, error) {
	var woofs []Woof
	err := database.DB.Where("post_id = ?", postID).Order("created_at asc").Find(&woofs).Error
	return woofs, err
}
