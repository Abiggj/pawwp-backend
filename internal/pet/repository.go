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
