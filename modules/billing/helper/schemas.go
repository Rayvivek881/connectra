package helper

import "time"

// BillingInfoResponse represents billing information for a user
type BillingInfoResponse struct {
	Credits              int        `json:"credits"`
	CreditsUsed          int        `json:"credits_used"`
	CreditsLimit         int        `json:"credits_limit"`
	SubscriptionPlan     string     `json:"subscription_plan"`
	SubscriptionPeriod   *string    `json:"subscription_period,omitempty"`
	SubscriptionStatus   string     `json:"subscription_status"`
	SubscriptionStartedAt *time.Time `json:"subscription_started_at,omitempty"`
	SubscriptionEndsAt   *time.Time `json:"subscription_ends_at,omitempty"`
	UsagePercentage      float64    `json:"usage_percentage"`
}

// Savings represents savings information for a period
type Savings struct {
	Amount     float64 `json:"amount"`
	Percentage int     `json:"percentage"`
}

// SubscriptionPeriodResponse represents a subscription period
type SubscriptionPeriodResponse struct {
	Period        string   `json:"period"`
	Credits       int      `json:"credits"`
	RatePerCredit float64  `json:"rate_per_credit"`
	Price         float64  `json:"price"`
	Savings       *Savings `json:"savings,omitempty"`
}

// SubscriptionPlanResponse represents a subscription plan with all periods
type SubscriptionPlanResponse struct {
	Tier     string                           `json:"tier"`
	Name     string                           `json:"name"`
	Category string                           `json:"category"`
	Periods  map[string]SubscriptionPeriodResponse `json:"periods"`
}

// SubscriptionPlansResponse represents the response for listing plans
type SubscriptionPlansResponse struct {
	Plans []SubscriptionPlanResponse `json:"plans"`
}

// AddonPackageResponse represents an addon package
type AddonPackageResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Credits       int     `json:"credits"`
	RatePerCredit float64 `json:"rate_per_credit"`
	Price         float64 `json:"price"`
}

// AddonPackagesResponse represents the response for listing addon packages
type AddonPackagesResponse struct {
	Packages []AddonPackageResponse `json:"packages"`
}

// SubscribeRequest represents a subscription request
type SubscribeRequest struct {
	Tier   string `json:"tier" binding:"required"`
	Period string `json:"period" binding:"required"`
}

// SubscribeResponse represents a subscription response
type SubscribeResponse struct {
	Message            string     `json:"message"`
	SubscriptionPlan   string     `json:"subscription_plan"`
	SubscriptionPeriod string     `json:"subscription_period"`
	Credits            int        `json:"credits"`
	SubscriptionEndsAt *time.Time `json:"subscription_ends_at,omitempty"`
}

// AddonPurchaseRequest represents an addon purchase request
type AddonPurchaseRequest struct {
	PackageID string `json:"package_id" binding:"required"`
}

// AddonPurchaseResponse represents an addon purchase response
type AddonPurchaseResponse struct {
	Message      string `json:"message"`
	Package      string `json:"package"`
	CreditsAdded int    `json:"credits_added"`
	TotalCredits int    `json:"total_credits"`
}

// CancelSubscriptionResponse represents a cancel subscription response
type CancelSubscriptionResponse struct {
	Message           string `json:"message"`
	SubscriptionStatus string `json:"subscription_status"`
}

// InvoiceItem represents an invoice item
type InvoiceItem struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Description *string   `json:"description,omitempty"`
}

// InvoiceListResponse represents the response for listing invoices
type InvoiceListResponse struct {
	Invoices []InvoiceItem `json:"invoices"`
	Total    int           `json:"total"`
}

// Admin CRUD Schemas

// SubscriptionPeriodCreate represents a request to create/update a subscription plan period
type SubscriptionPeriodCreate struct {
	Period            string   `json:"period" binding:"required"`
	Credits           int      `json:"credits" binding:"required,min=1"`
	RatePerCredit     float64  `json:"rate_per_credit" binding:"required,min=0"`
	Price             float64  `json:"price" binding:"required,min=0"`
	SavingsAmount     *float64 `json:"savings_amount,omitempty"`
	SavingsPercentage *int     `json:"savings_percentage,omitempty"`
}

// SubscriptionPlanCreate represents a request to create a subscription plan
type SubscriptionPlanCreate struct {
	Tier      string                      `json:"tier" binding:"required"`
	Name      string                      `json:"name" binding:"required"`
	Category  string                      `json:"category" binding:"required"`
	IsActive  *bool                       `json:"is_active,omitempty"`
	Periods   []SubscriptionPeriodCreate  `json:"periods" binding:"required"`
}

// SubscriptionPlanUpdate represents a request to update a subscription plan
type SubscriptionPlanUpdate struct {
	Name     *string `json:"name,omitempty"`
	Category *string `json:"category,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// AddonPackageCreate represents a request to create an addon package
type AddonPackageCreate struct {
	ID            string  `json:"id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Credits       int     `json:"credits" binding:"required,min=1"`
	RatePerCredit float64 `json:"rate_per_credit" binding:"required,min=0"`
	Price         float64 `json:"price" binding:"required,min=0"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

// AddonPackageUpdate represents a request to update an addon package
type AddonPackageUpdate struct {
	Name          *string  `json:"name,omitempty"`
	Credits       *int     `json:"credits,omitempty"`
	RatePerCredit *float64 `json:"rate_per_credit,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

