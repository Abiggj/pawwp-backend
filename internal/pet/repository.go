package pet

import (
	"github.com/Abiggj/pawwp/internal/database"
	"github.com/google/uuid"
)

type Repository interface {
	Create(pet *Pet) error
	FindByUsername(username string) (*Pet, error)
	FindByID(id string) (*Pet, error)
	FindByAccountID(accountID uuid.UUID) ([]Pet, error)
	Update(pet *Pet) error
	Delete(id string) error

	CreateMedicalRecord(record *MedicalRecord) error
	GetMedicalRecordsByPetID(petID uuid.UUID) ([]MedicalRecord, error)
	DeleteMedicalRecord(recordID uuid.UUID, accountID uuid.UUID) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(pet *Pet) error {
	return database.DB.Create(pet).Error
}

func (r *repository) FindByUsername(username string) (*Pet, error) {
	var pet Pet
	err := database.DB.Where("username = ?", username).First(&pet).Error
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

func (r *repository) FindByID(id string) (*Pet, error) {
	var pet Pet
	err := database.DB.Where("id = ?", id).First(&pet).Error
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

func (r *repository) FindByAccountID(accountID uuid.UUID) ([]Pet, error) {
	var pets []Pet
	err := database.DB.Where("account_id = ?", accountID).Find(&pets).Error
	return pets, err
}

func (r *repository) Update(pet *Pet) error {
	return database.DB.Save(pet).Error
}

func (r *repository) Delete(id string) error {
	return database.DB.Delete(&Pet{}, "id = ?", id).Error
}

func (r *repository) CreateMedicalRecord(record *MedicalRecord) error {
	return database.DB.Create(record).Error
}

func (r *repository) GetMedicalRecordsByPetID(petID uuid.UUID) ([]MedicalRecord, error) {
	var records []MedicalRecord
	err := database.DB.Where("pet_id = ?", petID).Order("created_at desc").Find(&records).Error
	return records, err
}

func (r *repository) DeleteMedicalRecord(recordID uuid.UUID, accountID uuid.UUID) error {
	return database.DB.Delete(&MedicalRecord{}, "id = ? AND account_id = ?", recordID, accountID).Error
}
