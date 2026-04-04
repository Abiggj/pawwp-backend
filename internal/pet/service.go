package pet

import (
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	Create(req CreatePetRequest, accountID uuid.UUID) (*Pet, error)
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

func (s *service) Create(req CreatePetRequest, accountID uuid.UUID) (*Pet, error) {

	// DATA STRUCTURE CONCEPT:
	// Username uniqueness behaves like a hash-set constraint.
	// DB enforces uniqueness, but we check first to avoid exception.

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
