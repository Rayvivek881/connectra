package helper

import (
	"database/sql"
	"time"
	"vivek-ray/models"
)

// GeolocationData represents IP geolocation data from frontend
type GeolocationData struct {
	IP            *string  `json:"ip,omitempty"`
	Continent     *string  `json:"continent,omitempty"`
	ContinentCode *string  `json:"continent_code,omitempty"`
	Country       *string  `json:"country,omitempty"`
	CountryCode   *string  `json:"country_code,omitempty"`
	Region        *string  `json:"region,omitempty"`
	RegionName    *string  `json:"region_name,omitempty"`
	City          *string  `json:"city,omitempty"`
	District      *string  `json:"district,omitempty"`
	Zip           *string  `json:"zip,omitempty"`
	Lat           *float64 `json:"lat,omitempty"`
	Lon           *float64 `json:"lon,omitempty"`
	Timezone      *string  `json:"timezone,omitempty"`
	Offset        *int     `json:"offset,omitempty"`
	Currency      *string  `json:"currency,omitempty"`
	ISP           *string  `json:"isp,omitempty"`
	Org           *string  `json:"org,omitempty"`
	ASName        *string  `json:"asname,omitempty"`
	Reverse       *string  `json:"reverse,omitempty"`
	Device        *string  `json:"device,omitempty"`
	Proxy         *bool    `json:"proxy,omitempty"`
	Hosting       *bool    `json:"hosting,omitempty"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name        string             `json:"name" binding:"required,max=255"`
	Email       string             `json:"email" binding:"required,email"`
	Password    string             `json:"password" binding:"required,min=8,max=72"`
	Geolocation *GeolocationData   `json:"geolocation,omitempty"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email       string             `json:"email" binding:"required,email"`
	Password    string             `json:"password" binding:"required"`
	Geolocation *GeolocationData   `json:"geolocation,omitempty"`
}

// UserResponse represents user information in responses
type UserResponse struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
}

// SessionUserResponse represents user information in session responses
type SessionUserResponse struct {
	UUID          string     `json:"uuid"`
	Email         string     `json:"email"`
	LastSignInAt  *time.Time `json:"last_sign_in_at,omitempty"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	WeeklyReports *bool `json:"weeklyReports,omitempty"`
	NewLeadAlerts *bool `json:"newLeadAlerts,omitempty"`
}

// ProfileResponse represents full user profile response
type ProfileResponse struct {
	UUID          string                  `json:"uuid"`
	Name          *string                 `json:"name,omitempty"`
	Email         string                  `json:"email"`
	Role          *string                 `json:"role,omitempty"`
	AvatarURL     *string                 `json:"avatar_url,omitempty"`
	IsActive      bool                    `json:"is_active"`
	JobTitle      *string                 `json:"job_title,omitempty"`
	Bio           *string                 `json:"bio,omitempty"`
	Timezone      *string                 `json:"timezone,omitempty"`
	Notifications *NotificationPreferences `json:"notifications,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     *time.Time              `json:"updated_at,omitempty"`
}

// ProfileUpdateRequest represents partial profile update
type ProfileUpdateRequest struct {
	Name          *string                 `json:"name,omitempty" binding:"omitempty,max=255"`
	JobTitle      *string                 `json:"job_title,omitempty" binding:"omitempty,max=255"`
	Bio           *string                 `json:"bio,omitempty"`
	Timezone      *string                 `json:"timezone,omitempty" binding:"omitempty,max=100"`
	AvatarURL     *string                 `json:"avatar_url,omitempty"`
	Notifications *NotificationPreferences `json:"notifications,omitempty"`
	Role          *string                 `json:"role,omitempty" binding:"omitempty,max=50"`
}

// TokenResponse represents token response
type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
	Message      string       `json:"message"`
}

// SessionResponse represents current session information
type SessionResponse struct {
	User SessionUserResponse `json:"user"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	RefreshToken *string `json:"refresh_token,omitempty"`
}

// LogoutResponse represents logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

// AvatarUploadResponse represents avatar upload response
type AvatarUploadResponse struct {
	AvatarURL string          `json:"avatar_url"`
	Profile   ProfileResponse `json:"profile"`
	Message   string          `json:"message"`
}

// UserListItem represents user list item (for Super Admin)
type UserListItem struct {
	UUID              string     `json:"uuid"`
	Email             string     `json:"email"`
	Name              *string    `json:"name,omitempty"`
	Role              *string    `json:"role,omitempty"`
	Credits           int        `json:"credits"`
	SubscriptionPlan  *string    `json:"subscription_plan,omitempty"`
	SubscriptionPeriod *string   `json:"subscription_period,omitempty"`
	SubscriptionStatus *string   `json:"subscription_status,omitempty"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	LastSignInAt      *time.Time `json:"last_sign_in_at,omitempty"`
}

// UserListResponse represents user list response
type UserListResponse struct {
	Users []UserListItem `json:"users"`
	Total int            `json:"total"`
}

// UpdateUserRoleRequest represents request to update user role
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserCreditsRequest represents request to update user credits
type UpdateUserCreditsRequest struct {
	Credits int `json:"credits" binding:"required,min=0"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers   int            `json:"total_users"`
	ActiveUsers  int            `json:"active_users"`
	UsersByRole  map[string]int `json:"users_by_role"`
	UsersByPlan  map[string]int `json:"users_by_plan"`
}

// UserHistoryItem represents a single user history record
type UserHistoryItem struct {
	ID            int       `json:"id"`
	UserID        string    `json:"user_id"`
	UserEmail     *string   `json:"user_email,omitempty"`
	UserName      *string   `json:"user_name,omitempty"`
	EventType     string    `json:"event_type"`
	IP            *string   `json:"ip,omitempty"`
	Continent     *string   `json:"continent,omitempty"`
	ContinentCode *string   `json:"continent_code,omitempty"`
	Country       *string   `json:"country,omitempty"`
	CountryCode   *string   `json:"country_code,omitempty"`
	Region        *string   `json:"region,omitempty"`
	RegionName    *string   `json:"region_name,omitempty"`
	City          *string   `json:"city,omitempty"`
	District      *string   `json:"district,omitempty"`
	Zip           *string   `json:"zip,omitempty"`
	Lat           *float64  `json:"lat,omitempty"`
	Lon           *float64  `json:"lon,omitempty"`
	Timezone      *string   `json:"timezone,omitempty"`
	Currency      *string   `json:"currency,omitempty"`
	ISP           *string   `json:"isp,omitempty"`
	Org           *string   `json:"org,omitempty"`
	ASName        *string   `json:"asname,omitempty"`
	Reverse       *string   `json:"reverse,omitempty"`
	Device        *string   `json:"device,omitempty"`
	Proxy         *bool     `json:"proxy,omitempty"`
	Hosting       *bool     `json:"hosting,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// UserHistoryListResponse represents paginated user history list response
type UserHistoryListResponse struct {
	Items  []UserHistoryItem `json:"items"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

// ToUserHistory converts GeolocationData to UserHistory model
func (g *GeolocationData) ToUserHistory(userID string, eventType models.UserHistoryEventType) *models.UserHistory {
	if g == nil {
		return nil
	}

	history := &models.UserHistory{
		UserID:        userID,
		EventType:     eventType,
		IP:            g.IP,
		Continent:     g.Continent,
		ContinentCode: g.ContinentCode,
		Country:       g.Country,
		CountryCode:   g.CountryCode,
		Region:        g.Region,
		RegionName:    g.RegionName,
		City:          g.City,
		District:      g.District,
		Zip:           g.Zip,
		Timezone:      g.Timezone,
		Currency:      g.Currency,
		ISP:           g.ISP,
		Org:           g.Org,
		ASName:        g.ASName,
		Reverse:       g.Reverse,
		Device:        g.Device,
		Proxy:         g.Proxy,
		Hosting:       g.Hosting,
	}

	// Convert lat/lon to sql.NullFloat64
	if g.Lat != nil {
		history.Lat = &sql.NullFloat64{
			Float64: *g.Lat,
			Valid:   true,
		}
	}
	if g.Lon != nil {
		history.Lon = &sql.NullFloat64{
			Float64: *g.Lon,
			Valid:   true,
		}
	}

	return history
}

