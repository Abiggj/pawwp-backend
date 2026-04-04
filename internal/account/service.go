package account

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req RegisterRequest) (*Account, error)
	Login(req LoginRequest) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type RegisterRequest struct {
	Type     string `json:"type"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (s *service) Register(req RegisterRequest) (*Account, error) {

	// DATA STRUCTURE CONCEPT:
	// Hashing function (bcrypt) uses cryptographic hashing
	// Similar conceptual use as hash maps for data transformation
	// We transform password into fixed-length hash (one-way function)

	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Type:         req.Type,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Bio:          req.Bio,
	}

	err = s.repo.Create(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *service) Login(req LoginRequest) (string, error) {

	account, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"account_id": account.ID.String(),
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
