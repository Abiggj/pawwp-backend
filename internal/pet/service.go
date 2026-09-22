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

	CreateMedicalRecord(req CreateMedicalRecordRequest, petID uuid.UUID, accountID uuid.UUID) (*MedicalRecord, error)
	GetMedicalRecords(petID uuid.UUID) ([]MedicalRecord, error)
	DeleteMedicalRecord(recordID uuid.UUID, accountID uuid.UUID) error
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

type CreateMedicalRecordRequest struct {
	RecordType       string `json:"record_type"`
	Title            string `json:"title"`
	Veterinarian     string `json:"veterinarian"`
	DateAdministered string `json:"date_administered"`
	ExpiryDate       string `json:"expiry_date"`
	Dosage           string `json:"dosage"`
	Notes            string `json:"notes"`
	DocumentURL      string `json:"document_url"`
	Status           string `json:"status"`
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

func (s *service) CreateMedicalRecord(req CreateMedicalRecordRequest, petID uuid.UUID, accountID uuid.UUID) (*MedicalRecord, error) {
	pet, err := s.repo.FindByID(petID.String())
	if err != nil {
		return nil, errors.New("pet not found")
	}

	if pet.AccountID != accountID {
		return nil, errors.New("unauthorized: only pet owner can add medical records")
	}

	if req.Title == "" {
		return nil, errors.New("record title is required")
	}
	if req.RecordType == "" {
		req.RecordType = "vaccination"
	}
	if req.Status == "" {
		req.Status = "active"
	}

	record := &MedicalRecord{
		PetID:            petID,
		AccountID:        accountID,
		RecordType:       req.RecordType,
		Title:            req.Title,
		Veterinarian:     req.Veterinarian,
		DateAdministered: req.DateAdministered,
		ExpiryDate:       req.ExpiryDate,
		Dosage:           req.Dosage,
		Notes:            req.Notes,
		DocumentURL:      req.DocumentURL,
		Status:           req.Status,
	}

	err = s.repo.CreateMedicalRecord(record)
	return record, err
}

func (s *service) GetMedicalRecords(petID uuid.UUID) ([]MedicalRecord, error) {
	return s.repo.GetMedicalRecordsByPetID(petID)
}

func (s *service) DeleteMedicalRecord(recordID uuid.UUID, accountID uuid.UUID) error {
	return s.repo.DeleteMedicalRecord(recordID, accountID)
}
