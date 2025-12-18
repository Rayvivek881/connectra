package service

import (
	"context"
	"errors"
	"time"
	"vivek-ray/conf"
	"vivek-ray/models"
	"vivek-ray/modules/auth/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repository.NewAuthRepository(),
	}
}

type RegisterRequest struct {
	Email       string             `json:"email"`
	Password    string             `json:"password"`
	Role        models.Role        `json:"role"`
	Geolocation models.Geolocation `json:"geolocation"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil && existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return nil, err
	}

	user := &models.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		Role:        req.Role, // Default role
		Credits:     50,       // Default credits
		Geolocation: req.Geolocation,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateToken(user, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	// Parse token to get expiration time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration time")
	}

	expiresAt := time.Unix(int64(exp), 0)
	return s.repo.AddToBlacklist(ctx, tokenString, expiresAt)
}

func (s *AuthService) generateToken(user *models.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.AppConfig.JWTSecret))
}
