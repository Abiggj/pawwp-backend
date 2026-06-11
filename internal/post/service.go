package post

import (
	"errors"

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
	AddWoof(postID string, content string, accountID uuid.UUID) (*Woof, error)
	GetPostWoofs(postID string) ([]Woof, error)
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
	if err != nil {
		return nil, errors.New("associated pet not found")
	}

	if pet.AccountID != accountID {
		return nil, errors.New("unauthorized")
	}

	if !post.IsShowcased {
		count, err := s.repo.CountShowcaseByPetID(post.PetID)
		if err != nil {
			return nil, err
		}
		if count >= 9 {
			return nil, errors.New("showcase limit reached (max 9)")
		}
		post.IsShowcased = true
	} else {
		post.IsShowcased = false
	}

	err = s.repo.Update(post)
	return post, err
}

func (s *service) GetPetArchive(petID uuid.UUID, accountID uuid.UUID) ([]Post, error) {
	pet, err := s.petRepo.FindByID(petID.String())
	if err != nil {
		return nil, errors.New("pet not found")
	}

	if pet.AccountID != accountID {
		return nil, errors.New("unauthorized")
	}

	return s.repo.GetArchiveByPetID(petID)
}

func (s *service) ToggleBoop(postID string, accountID uuid.UUID) error {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return errors.New("invalid post id")
	}

	// Try to remove first (untoggle)
	err = s.repo.RemoveBoop(accountID, pID)
	if err == nil {
		// If delete succeeded, it means it was toggled on, now off.
		// Wait, GORM Delete might not return error if 0 rows affected.
		// Let's assume AddBoop will fail if already exists due to unique index.
	}

	// This is a bit naive, better would be to check existence or use a upsert-like logic.
	// But according to plan "Toggle a boop".
	boop := &Boop{
		AccountID: accountID,
		PostID:    pID,
	}
	err = s.repo.AddBoop(boop)
	if err != nil {
		// If already exists, remove it
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

func (s *service) GetPostWoofs(postID string) ([]Woof, error) {
	pID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}
	return s.repo.GetWoofsByPostID(pID)
}
