package service

import (
	"context"
	"errors"
	"time"
	"vivek-ray/conf"
	"vivek-ray/models"
	authrepo "vivek-ray/modules/auth/repository"
	"vivek-ray/modules/users/helper"
	userrepo "vivek-ray/modules/users/repository"
	userservice "vivek-ray/modules/users/service"
	"vivek-ray/utilities"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         *authrepo.AuthRepository
	userRepo     *userrepo.UserRepository
	profileRepo  *userrepo.UserProfileRepository
	historySvc   *userservice.HistoryService
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo:        authrepo.NewAuthRepository(),
		userRepo:    userrepo.NewUserRepository(),
		profileRepo: userrepo.NewUserProfileRepository(),
		historySvc:  userservice.NewHistoryService(),
	}
}

// Register registers a new user and creates their profile
func (s *AuthService) Register(ctx context.Context, req *helper.RegisterRequest) (*models.User, string, string, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, "", "", errors.New("email already exists")
	}

	// Validate password
	if len(req.Password) < 8 {
		return nil, "", "", errors.New("password must be at least 8 characters")
	}
	if len(req.Password) > 72 {
		return nil, "", "", errors.New("password must be at most 72 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return nil, "", "", err
	}

	// Generate UUIDs
	userID := utilities.GenerateUUID()
	userUUID := utilities.GenerateUUID()

	// Create user
	user := &models.User{
		ID:             userID,
		UUID:           userUUID,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
		Name:           &req.Name,
		IsActive:       true,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, "", "", err
	}

	// Create default profile
	defaultRole := "FreeUser"
	defaultPlan := "free"
	defaultStatus := "active"
	profile := &models.UserProfile{
		UserID:             userUUID,
		Role:               &defaultRole,
		Credits:            50,
		SubscriptionPlan:   &defaultPlan,
		SubscriptionStatus: &defaultStatus,
		Notifications: &models.NotificationPreferences{
			WeeklyReports: boolPtr(true),
			NewLeadAlerts: boolPtr(true),
		},
	}

	if err := s.profileRepo.CreateProfile(ctx, profile); err != nil {
		return nil, "", "", err
	}

	// Generate tokens
	accessToken, err := s.generateToken(userUUID, 30*time.Minute, "access")
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateToken(userUUID, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, "", "", err
	}

	// Record registration history
	if req.Geolocation != nil {
		if err := s.historySvc.RecordRegistration(ctx, userUUID, req.Geolocation); err != nil {
			log.Warn().Err(err).Msg("Failed to record registration history")
		}
	}

	return user, accessToken, refreshToken, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *helper.LoginRequest) (*models.User, string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}
	if user == nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, "", "", errors.New("user account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	// Update last sign in
	now := time.Now()
	user.LastSignInAt = &now
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		log.Warn().Err(err).Msg("Failed to update last sign in")
	}

	// Generate tokens
	accessToken, err := s.generateToken(user.UUID, 30*time.Minute, "access")
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateToken(user.UUID, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, "", "", err
	}

	// Record login history
	if req.Geolocation != nil {
		if err := s.historySvc.RecordLogin(ctx, user.UUID, req.Geolocation); err != nil {
			log.Warn().Err(err).Msg("Failed to record login history")
		}
	}

	return user, accessToken, refreshToken, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Check if token is blacklisted
	isBlacklisted, err := s.repo.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}
	if isBlacklisted {
		return "", "", errors.New("token is invalid or expired")
	}

	// Parse and validate token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(conf.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", errors.New("invalid refresh token")
	}

	// Get user UUID
	userUUID, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// Verify user exists and is active
	user, err := s.userRepo.GetByUUID(ctx, userUUID)
	if err != nil || user == nil || !user.IsActive {
		return "", "", errors.New("token is invalid or expired")
	}

	// Generate new tokens (token rotation)
	newAccessToken, err := s.generateToken(userUUID, 30*time.Minute, "access")
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateToken(userUUID, 7*24*time.Hour, "refresh")
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout logs out a user and blacklists refresh token
func (s *AuthService) Logout(ctx context.Context, tokenString string, userID *string) error {
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
	return s.repo.AddToBlacklist(ctx, tokenString, expiresAt, userID)
}

// GetSession returns current session information
func (s *AuthService) GetSession(ctx context.Context, userUUID string) (*helper.SessionResponse, error) {
	user, err := s.userRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &helper.SessionResponse{
		User: helper.SessionUserResponse{
			UUID:         user.UUID,
			Email:        user.Email,
			LastSignInAt: user.LastSignInAt,
		},
	}, nil
}

// generateToken generates a JWT token
func (s *AuthService) generateToken(userUUID string, duration time.Duration, tokenType string) (string, error) {
	expireMinutes := conf.AppConfig.AccessTokenExpire
	if expireMinutes == 0 {
		expireMinutes = 30
	}
	if tokenType == "refresh" {
		expireDays := conf.AppConfig.RefreshTokenExpire
		if expireDays == 0 {
			expireDays = 7
		}
		duration = time.Duration(expireDays) * 24 * time.Hour
	} else {
		duration = time.Duration(expireMinutes) * time.Minute
	}

	claims := jwt.MapClaims{
		"sub":   userUUID,
		"type":  tokenType,
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.AppConfig.JWTSecret))
}

func boolPtr(b bool) *bool {
	return &b
}
