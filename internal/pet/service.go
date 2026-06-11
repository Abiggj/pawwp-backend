package pet

import (
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	Create(req CreatePetRequest, accountID uuid.UUID) (*Pet, error)
	GetByID(id string) (*Pet, error)
	GetByAccountID(accountID uuid.UUID) ([]Pet, error)
	Update(id string, req UpdatePetRequest, accountID uuid.UUID) (*Pet, error)
	Delete(id string, accountID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type CreatePetRequest struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	Species        string `json:"species"`
	Breed          string `json:"breed"`
	Bio            string `json:"bio"`
	AdoptionStatus string `json:"adoption_status"`
}

type UpdatePetRequest struct {
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	AdoptionStatus string `json:"adoption_status"`
}

func (s *service) Create(req CreatePetRequest, accountID uuid.UUID) (*Pet, error) {
	existing, _ := s.repo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already taken")
	}

	pet := &Pet{
		AccountID:      accountID,
		Username:       req.Username,
		Name:           req.Name,
		Species:        req.Species,
		Breed:          req.Breed,
		Bio:            req.Bio,
		AdoptionStatus: req.AdoptionStatus,
	}

	err := s.repo.Create(pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *service) GetByID(id string) (*Pet, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetByAccountID(accountID uuid.UUID) ([]Pet, error) {
	return s.repo.FindByAccountID(accountID)
}

func (s *service) Update(id string, req UpdatePetRequest, accountID uuid.UUID) (*Pet, error) {
	pet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if pet.AccountID != accountID {
		return nil, errors.New("unauthorized")
	}

	pet.Name = req.Name
	pet.Bio = req.Bio
	pet.AdoptionStatus = req.AdoptionStatus

	err = s.repo.Update(pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *service) Delete(id string, accountID uuid.UUID) error {
	pet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if pet.AccountID != accountID {
		return errors.New("unauthorized")
	}

	return s.repo.Delete(id)
}
