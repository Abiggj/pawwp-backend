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
	Login(req LoginRequest) (*AuthResponse, error)
	Refresh(refreshToken string) (*AuthResponse, error)
	GetByID(id string) (*Account, error)
	Update(id string, req UpdateRequest) (*Account, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Type     string `json:"type"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UpdateRequest struct {
	Bio string `json:"bio"`
}

func (s *service) Register(req RegisterRequest) (*Account, error) {
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

func (s *service) Login(req LoginRequest) (*AuthResponse, error) {
	account, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAndSaveTokens(account)
}

func (s *service) Refresh(refreshToken string) (*AuthResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	accountID := claims["account_id"].(string)
	account, err := s.repo.FindByID(accountID)
	if err != nil || account.RefreshToken != refreshToken {
		return nil, errors.New("invalid refresh token")
	}

	return s.generateAndSaveTokens(account)
}

func (s *service) GetByID(id string) (*Account, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(id string, req UpdateRequest) (*Account, error) {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	account.Bio = req.Bio
	err = s.repo.Update(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *service) generateAndSaveTokens(account *Account) (*AuthResponse, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Access Token
	accessClaims := jwt.MapClaims{
		"account_id": account.ID.String(),
		"type":       "access",
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"account_id": account.ID.String(),
		"type":       "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	account.RefreshToken = signedRefresh
	err = s.repo.Update(account)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  signedAccess,
		RefreshToken: signedRefresh,
	}, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
