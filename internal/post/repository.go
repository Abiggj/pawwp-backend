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
	GetBoopsByPostID(postID uuid.UUID) ([]BoopDetail, error)
	AddWoof(woof *Woof) error
	GetWoofsByPostID(postID uuid.UUID) ([]Woof, error)
	GetWoofDetailsByPostID(postID uuid.UUID) ([]WoofDetail, error)
	RecordAuthenticView(postID uuid.UUID, accountID *uuid.UUID, viewerKey string, dwellSeconds int) (int64, bool, error)
	GetPostViewsCount(postID uuid.UUID) (int64, error)
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

func (r *repository) GetBoopsByPostID(postID uuid.UUID) ([]BoopDetail, error) {
	var boops []BoopDetail
	err := database.DB.Table("boops").
		Select("boops.id, boops.post_id, boops.account_id, accounts.email as user_name, '' as user_avatar, accounts.type as user_type, boops.created_at as booped_at").
		Joins("JOIN accounts ON accounts.id = boops.account_id").
		Where("boops.post_id = ?", postID).
		Order("boops.created_at desc").
		Scan(&boops).Error
	return boops, err
}

func (r *repository) AddWoof(woof *Woof) error {
	return database.DB.Create(woof).Error
}

func (r *repository) GetWoofsByPostID(postID uuid.UUID) ([]Woof, error) {
	var woofs []Woof
	err := database.DB.Where("post_id = ?", postID).Order("created_at asc").Find(&woofs).Error
	return woofs, err
}

func (r *repository) GetWoofDetailsByPostID(postID uuid.UUID) ([]WoofDetail, error) {
	var woofs []WoofDetail
	err := database.DB.Table("woofs").
		Select("woofs.id, woofs.post_id, woofs.account_id, accounts.email as author_name, '' as author_avatar, accounts.type as author_type, woofs.content, woofs.created_at").
		Joins("JOIN accounts ON accounts.id = woofs.account_id").
		Where("woofs.post_id = ?", postID).
		Order("woofs.created_at asc").
		Scan(&woofs).Error
	return woofs, err
}

func (r *repository) RecordAuthenticView(postID uuid.UUID, accountID *uuid.UUID, viewerKey string, dwellSeconds int) (int64, bool, error) {
	if dwellSeconds < 2 {
		count, _ := r.GetPostViewsCount(postID)
		return count, false, nil
	}

	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
	var existingCount int64
	database.DB.Model(&PostView{}).
		Where("post_id = ? AND viewer_key = ? AND created_at >= ?", postID, viewerKey, twentyFourHoursAgo).
		Count(&existingCount)

	if existingCount > 0 {
		count, _ := r.GetPostViewsCount(postID)
		return count, false, nil
	}

	view := &PostView{
		PostID:       postID,
		AccountID:    accountID,
		ViewerKey:    viewerKey,
		DwellSeconds: dwellSeconds,
		IsAuthentic:  true,
		CreatedAt:    time.Now(),
	}
	if err := database.DB.Create(view).Error; err != nil {
		count, _ := r.GetPostViewsCount(postID)
		return count, false, err
	}

	count, err := r.GetPostViewsCount(postID)
	return count, true, err
}

func (r *repository) GetPostViewsCount(postID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&PostView{}).Where("post_id = ? AND is_authentic = ?", postID, true).Count(&count).Error
	return count, err
}
