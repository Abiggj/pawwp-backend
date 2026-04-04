package pet

import (
	"github.com/Abiggj/pawwp/internal/database"
)

type Repository interface {
	Create(pet *Pet) error
	FindByUsername(username string) (*Pet, error)
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
