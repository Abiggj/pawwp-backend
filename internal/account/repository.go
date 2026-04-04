package account

import (
	"github.com/Abiggj/pawwp/internal/database"
)

type Repository interface {
	Create(account *Account) error
	FindByEmail(email string) (*Account, error)
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
