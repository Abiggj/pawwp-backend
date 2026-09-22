package post

import (
	"errors"
	"time"

	"github.com/Abiggj/pawwp/internal/pet"
	"github.com/google/uuid"
)

type Service interface {
	CreatePost(req CreatePostRequest, accountID uuid.UUID) (*Post, error)
	Get48hFeed() ([]Post, error)
	GetPetShowcase(petID uuid.UUID) ([]Post, error)
	ToggleShowcase(postID string, accountID uuid.UUID) (*Post, error)
	GetPetArchive(petID uuid.UUID, accountID uuid.UUID) ([]Post, error)

	ToggleBoop(postID string, accountID uuid.UUID) error
	GetPostBoops(postID string) ([]BoopDetail, error)
	AddWoof(postID string, content string, accountID uuid.UUID) (*Woof, error)
	GetPostWoofs(postID string) ([]WoofDetail, error)
	RecordView(postID string, accountID *uuid.UUID, viewerKey string, dwellSeconds int) (int64, bool, error)
	GetPostViews(postID string) (int64, error)
}

type service struct {
	repo    Repository
	petRepo pet.Repository
}

func NewService(repo Repository, petRepo pet.Repository) Service {
	return &service{repo: repo, petRepo: petRepo}
}

type CreatePostRequest struct {
	PetID    uuid.UUID `json:"pet_id"`
	MediaURL string    `json:"media_url"`
	Caption  string    `json:"caption"`
}

func (s *service) CreatePost(req CreatePostRequest, accountID uuid.UUID) (*Post, error) {
	pet, err := s.petRepo.FindByID(req.PetID.String())
	if err != nil {
		return nil, errors.New("pet not found")
	}

	if pet.AccountID != accountID {
		return nil, errors.New("unauthorized: you do not own this pet")
	}

	post := &Post{
		PetID:    req.PetID,
		MediaURL: req.MediaURL,
		Caption:  req.Caption,
	}

	err = s.repo.Create(post)
	return post, err
}

func (s *service) Get48hFeed() ([]Post, error) {
	return s.repo.Get48hFeed()
}

func (s *service) GetPetShowcase(petID uuid.UUID) ([]Post, error) {
	return s.repo.GetShowcaseByPetID(petID)
}

func (s *service) ToggleShowcase(postID string, accountID uuid.UUID) (*Post, error) {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	pet, err := s.petRepo.FindByID(post.PetID.String())
	if err != nil || pet.AccountID != accountID {
		return nil, errors.New("unauthorized: you do not own this post")
	}

	if !post.IsShowcased {
		count, err := s.repo.CountShowcaseByPetID(post.PetID)
		if err != nil {
			return nil, err
		}
		if count >= 9 {
			return nil, errors.New("showcase is full (max 9 posts)")
		}
	}

	post.IsShowcased = !post.IsShowcased
	err = s.repo.Update(post)
	return post, err
}

func (s *service) GetPetArchive(petID uuid.UUID, accountID uuid.UUID) ([]Post, error) {
	pet, err := s.petRepo.FindByID(petID.String())
	if err != nil || pet.AccountID != accountID {
		return nil, errors.New("unauthorized: you do not own this pet")
	}

	return s.repo.GetArchiveByPetID(petID)
}

func (s *service) ToggleBoop(postID string, accountID uuid.UUID) error {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return errors.New("invalid post id")
	}

	post, err := s.repo.FindByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	fortyEightHoursAgo := time.Now().Add(-48 * time.Hour)
	if post.CreatedAt.Before(fortyEightHoursAgo) {
		return errors.New("this post has expired (older than 48 hours)")
	}

	boop := &Boop{
		AccountID: accountID,
		PostID:    pID,
	}

	err = s.repo.AddBoop(boop)
	if err != nil {
		return s.repo.RemoveBoop(accountID, pID)
	}

	return nil
}

func (s *service) AddWoof(postID string, content string, accountID uuid.UUID) (*Woof, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}

	woof := &Woof{
		AccountID: accountID,
		PostID:    pID,
		Content:   content,
	}

	err = s.repo.AddWoof(woof)
	return woof, err
}

func (s *service) GetPostBoops(postID string) ([]BoopDetail, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}
	return s.repo.GetBoopsByPostID(pID)
}

func (s *service) GetPostWoofs(postID string) ([]WoofDetail, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}
	return s.repo.GetWoofDetailsByPostID(pID)
}

func (s *service) RecordView(postID string, accountID *uuid.UUID, viewerKey string, dwellSeconds int) (int64, bool, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return 0, false, errors.New("invalid post id")
	}
	if viewerKey == "" {
		if accountID != nil {
			viewerKey = accountID.String()
		} else {
			viewerKey = "guest"
		}
	}
	return s.repo.RecordAuthenticView(pID, accountID, viewerKey, dwellSeconds)
}

func (s *service) GetPostViews(postID string) (int64, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return 0, errors.New("invalid post id")
	}
	return s.repo.GetPostViewsCount(pID)
}
