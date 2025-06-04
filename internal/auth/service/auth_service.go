package service

import (
	"time"

	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/model"
	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(req *model.SignupRequest) (*model.User, error)
	Login(req *model.LoginRequest) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, tokenTTL time.Duration) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

func (s *authService) SignUp(req *model.SignupRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Role:         req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(req *model.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"role": user.Role,
		"exp":  time.Now().Add(s.tokenTTL).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
