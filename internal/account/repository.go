package account

import (
	"github.com/Abiggj/pawwp/internal/database"
)

type Repository interface {
	Create(account *Account) error
	FindByEmail(email string) (*Account, error)
	FindByID(id string) (*Account, error)
	Update(account *Account) error
	Delete(id string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(account *Account) error {
	return database.DB.Create(account).Error
}

func (r *repository) FindByEmail(email string) (*Account, error) {
	var account Account
	err := database.DB.Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) FindByID(id string) (*Account, error) {
	var account Account
	err := database.DB.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) Update(account *Account) error {
	return database.DB.Save(account).Error
}

func (r *repository) Delete(id string) error {
	return database.DB.Delete(&Account{}, "id = ?", id).Error
}
