package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nadduli/go-modules/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type Service struct {
	repo UserRepository

	jwtSecret string
}

func NewService(repo UserRepository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) RegisterUser(ctx context.Context, username, email, password string) (models.PublicUser, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.PublicUser{}, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedBytes),
		Role:     "user",
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return models.PublicUser{}, err
	}

	return models.ToPublic(*user), nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, models.PublicUser, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", models.PublicUser{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", models.PublicUser{}, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return "", models.PublicUser{}, err
	}

	return token, models.ToPublic(*user), nil
}

func (s *Service) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
