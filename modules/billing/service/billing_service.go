package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vivek-ray/models"
	"vivek-ray/modules/billing/helper"
	"vivek-ray/modules/users/repository"

	"github.com/rs/zerolog/log"
)

const (
	MONTHLY   = "monthly"
	QUARTERLY = "quarterly"
	YEARLY    = "yearly"
)

// BillingService handles billing and subscription business logic
type BillingService struct {
	profileRepo        *repository.UserProfileRepository
	planRepo           *models.SubscriptionPlanRepository
	periodRepo         *models.SubscriptionPlanPeriodRepository
	addonRepo          *models.AddonPackageRepository
}

// NewBillingService creates a new billing service
func NewBillingService() *BillingService {
	return &BillingService{
		profileRepo: repository.NewUserProfileRepository(),
		planRepo:    models.NewSubscriptionPlanRepository(),
		periodRepo:  models.NewSubscriptionPlanPeriodRepository(),
		addonRepo:   models.NewAddonPackageRepository(),
	}
}

// calculateSubscriptionEndDate calculates subscription end date based on billing period
func calculateSubscriptionEndDate(period string, startDate time.Time) time.Time {
	switch period {
	case MONTHLY:
		return startDate.AddDate(0, 1, 0) // 30 days approximation
	case QUARTERLY:
		return startDate.AddDate(0, 3, 0) // 90 days
	case YEARLY:
		return startDate.AddDate(1, 0, 0) // 365 days
	default:
		return startDate.AddDate(0, 1, 0) // Default to monthly
	}
}

// GetBillingInfo retrieves billing information for a user
func (s *BillingService) GetBillingInfo(ctx context.Context, userID string) (*helper.BillingInfoResponse, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("user profile not found")
	}

	subscriptionTier := "free"
	if profile.SubscriptionPlan != nil {
		subscriptionTier = *profile.SubscriptionPlan
	}

	subscriptionPeriod := MONTHLY
	if profile.SubscriptionPeriod != nil {
		subscriptionPeriod = *profile.SubscriptionPeriod
	}

	// Calculate credits limit based on subscription
	var creditsLimit int
	if subscriptionTier == "free" {
		creditsLimit = 50 // Initial free credits
	} else {
		// Try to get from database
		plan, err := s.planRepo.GetByTier(ctx, subscriptionTier)
		if err == nil && plan != nil && plan.IsActive {
			periodObj, err := s.periodRepo.GetByPlanAndPeriod(ctx, subscriptionTier, subscriptionPeriod)
			if err == nil && periodObj != nil {
				creditsLimit = periodObj.Credits
			} else {
				// Fallback to monthly if period not found
				periodObj, err := s.periodRepo.GetByPlanAndPeriod(ctx, subscriptionTier, MONTHLY)
				if err == nil && periodObj != nil {
					creditsLimit = periodObj.Credits
				} else {
					creditsLimit = 50 // Default to free tier
				}
			}
		} else {
			creditsLimit = 50 // Default to free tier
		}
	}

	currentCredits := profile.Credits
	if currentCredits < 0 {
		currentCredits = 0
	}

	creditsUsed := creditsLimit - currentCredits
	if creditsUsed < 0 {
		creditsUsed = 0
	}
	if creditsUsed > creditsLimit {
		creditsUsed = creditsLimit
	}

	var usagePercentage float64
	if creditsLimit > 0 {
		usagePercentage = float64(creditsUsed) / float64(creditsLimit) * 100
	}

	subscriptionStatus := "active"
	if profile.SubscriptionStatus != nil {
		subscriptionStatus = *profile.SubscriptionStatus
	}

	return &helper.BillingInfoResponse{
		Credits:              currentCredits,
		CreditsUsed:          creditsUsed,
		CreditsLimit:         creditsLimit,
		SubscriptionPlan:     subscriptionTier,
		SubscriptionPeriod:   profile.SubscriptionPeriod,
		SubscriptionStatus:   subscriptionStatus,
		SubscriptionStartedAt: profile.SubscriptionStartedAt,
		SubscriptionEndsAt:   profile.SubscriptionEndsAt,
		UsagePercentage:     usagePercentage,
	}, nil
}

// GetSubscriptionPlans retrieves all available subscription plans with periods
func (s *BillingService) GetSubscriptionPlans(ctx context.Context) (*helper.SubscriptionPlansResponse, error) {
	plans, err := s.planRepo.ListAllWithPeriods(ctx, false)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving subscription plans")
		return nil, err
	}

	response := &helper.SubscriptionPlansResponse{
		Plans: make([]helper.SubscriptionPlanResponse, 0, len(plans)),
	}

	for _, plan := range plans {
		planResponse := helper.SubscriptionPlanResponse{
			Tier:     plan.Tier,
			Name:     plan.Name,
			Category: plan.Category,
			Periods:  make(map[string]helper.SubscriptionPeriodResponse),
		}

		for _, period := range plan.Periods {
			var savings *helper.Savings
			if period.SavingsAmount != nil || period.SavingsPercentage != nil {
				savings = &helper.Savings{}
				if period.SavingsAmount != nil {
					savings.Amount = *period.SavingsAmount
				}
				if period.SavingsPercentage != nil {
					savings.Percentage = *period.SavingsPercentage
				}
			}

			planResponse.Periods[period.Period] = helper.SubscriptionPeriodResponse{
				Period:        period.Period,
				Credits:       period.Credits,
				RatePerCredit: period.RatePerCredit,
				Price:         period.Price,
				Savings:       savings,
			}
		}

		response.Plans = append(response.Plans, planResponse)
	}

	return response, nil
}

// GetAddonPackages retrieves all available addon packages
func (s *BillingService) GetAddonPackages(ctx context.Context) (*helper.AddonPackagesResponse, error) {
	packages, err := s.addonRepo.ListAll(ctx, false)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving addon packages")
		return nil, err
	}

	response := &helper.AddonPackagesResponse{
		Packages: make([]helper.AddonPackageResponse, 0, len(packages)),
	}

	for _, pkg := range packages {
		response.Packages = append(response.Packages, helper.AddonPackageResponse{
			ID:            pkg.ID,
			Name:          pkg.Name,
			Credits:       pkg.Credits,
			RatePerCredit: pkg.RatePerCredit,
			Price:         pkg.Price,
		})
	}

	return response, nil
}

// SubscribeToPlan subscribes a user to a subscription plan
func (s *BillingService) SubscribeToPlan(ctx context.Context, userID, tier, period string) (*helper.SubscribeResponse, error) {
	// Validate period
	if period != MONTHLY && period != QUARTERLY && period != YEARLY {
		return nil, fmt.Errorf("invalid period: %s. Must be one of: monthly, quarterly, yearly", period)
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("user profile not found")
	}

	// Get plan and period from database
	plan, err := s.planRepo.GetByTier(ctx, tier)
	if err != nil {
		return nil, err
	}
	if plan == nil || !plan.IsActive {
		return nil, fmt.Errorf("invalid tier: %s", tier)
	}

	periodObj, err := s.periodRepo.GetByPlanAndPeriod(ctx, tier, period)
	if err != nil {
		return nil, err
	}
	if periodObj == nil {
		return nil, fmt.Errorf("invalid period: %s for tier %s", period, tier)
	}

	// Update subscription
	now := time.Now()
	profile.SubscriptionPlan = &tier
	profile.SubscriptionPeriod = &period
	status := "active"
	profile.SubscriptionStatus = &status
	profile.SubscriptionStartedAt = &now
	endDate := calculateSubscriptionEndDate(period, now)
	profile.SubscriptionEndsAt = &endDate
	profile.Credits = periodObj.Credits // Set credits to plan amount

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return &helper.SubscribeResponse{
		Message:            fmt.Sprintf("Successfully subscribed to %s (%s)", plan.Name, period),
		SubscriptionPlan:   tier,
		SubscriptionPeriod: period,
		Credits:            periodObj.Credits,
		SubscriptionEndsAt: &endDate,
	}, nil
}

// PurchaseAddonCredits purchases addon credits for a user
func (s *BillingService) PurchaseAddonCredits(ctx context.Context, userID, packageID string) (*helper.AddonPurchaseResponse, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("user profile not found")
	}

	// Get addon package
	pkg, err := s.addonRepo.GetByID(ctx, packageID)
	if err != nil {
		return nil, err
	}
	if pkg == nil || !pkg.IsActive {
		return nil, fmt.Errorf("invalid package ID: %s", packageID)
	}

	// Add credits to existing balance
	currentCredits := profile.Credits
	profile.Credits = currentCredits + pkg.Credits

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return &helper.AddonPurchaseResponse{
		Message:      fmt.Sprintf("Successfully purchased %d addon credits.", pkg.Credits),
		Package:      packageID,
		CreditsAdded:  pkg.Credits,
		TotalCredits:  profile.Credits,
	}, nil
}

// CancelSubscription cancels a user's subscription
func (s *BillingService) CancelSubscription(ctx context.Context, userID string) (*helper.CancelSubscriptionResponse, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("user profile not found")
	}

	if profile.SubscriptionStatus != nil && *profile.SubscriptionStatus == "cancelled" {
		return nil, errors.New("subscription is already cancelled")
	}

	// Mark as cancelled but keep active until end date
	status := "cancelled"
	profile.SubscriptionStatus = &status

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return &helper.CancelSubscriptionResponse{
		Message:            "Subscription cancelled. You will retain access until the end of your billing period.",
		SubscriptionStatus: "cancelled",
	}, nil
}

// GetInvoices retrieves invoice history for a user
func (s *BillingService) GetInvoices(ctx context.Context, userID string, limit, offset int) (*helper.InvoiceListResponse, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("user profile not found")
	}

	// TODO: In production, fetch invoices from payment processor (e.g., Stripe)
	// This is a simplified mock implementation
	invoices := make([]helper.InvoiceItem, 0)

	if profile.SubscriptionStartedAt != nil && profile.SubscriptionPlan != nil && *profile.SubscriptionPlan != "free" {
		tier := *profile.SubscriptionPlan
		period := MONTHLY
		if profile.SubscriptionPeriod != nil {
			period = *profile.SubscriptionPeriod
		}

		periodObj, err := s.periodRepo.GetByPlanAndPeriod(ctx, tier, period)
		if err == nil && periodObj != nil {
			plan, err := s.planRepo.GetByTier(ctx, tier)
			if err == nil && plan != nil {
				planPrice := periodObj.Price
				planName := plan.Name

				subscriptionStart := *profile.SubscriptionStartedAt
				now := time.Now()
				endDate := now
				if profile.SubscriptionEndsAt != nil && profile.SubscriptionEndsAt.Before(now) {
					endDate = *profile.SubscriptionEndsAt
				}

				// Generate invoices based on period
				currentDate := subscriptionStart
				invoiceNumber := 1
				var daysPerPeriod int
				switch period {
				case MONTHLY:
					daysPerPeriod = 30
				case QUARTERLY:
					daysPerPeriod = 90
				case YEARLY:
					daysPerPeriod = 365
				default:
					daysPerPeriod = 30
				}

				for currentDate.Before(endDate) && invoiceNumber <= 12 {
					var invoiceStatus string
					if profile.SubscriptionStatus != nil && *profile.SubscriptionStatus == "cancelled" && profile.SubscriptionEndsAt != nil && currentDate.After(*profile.SubscriptionEndsAt) {
						invoiceStatus = "failed"
					} else if currentDate.Before(now) || currentDate.Equal(now) {
						invoiceStatus = "paid"
					} else {
						invoiceStatus = "pending"
					}

					description := fmt.Sprintf("Subscription to %s (%s)", planName, period)
					invoices = append(invoices, helper.InvoiceItem{
						ID:          fmt.Sprintf("inv_%s_%03d", userID[:8], invoiceNumber),
						Amount:      planPrice,
						Status:      invoiceStatus,
						CreatedAt:   currentDate,
						Description: &description,
					})

					currentDate = currentDate.AddDate(0, 0, daysPerPeriod)
					invoiceNumber++
				}
			}
		}
	}

	// Sort invoices by date (newest first)
	for i := 0; i < len(invoices)-1; i++ {
		for j := i + 1; j < len(invoices); j++ {
			if invoices[i].CreatedAt.Before(invoices[j].CreatedAt) {
				invoices[i], invoices[j] = invoices[j], invoices[i]
			}
		}
	}

	// Apply pagination
	total := len(invoices)
	start := offset
	if start > total {
		start = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	paginatedInvoices := invoices[start:end]

	return &helper.InvoiceListResponse{
		Invoices: paginatedInvoices,
		Total:    total,
	}, nil
}

// Admin CRUD methods for Subscription Plans

// CreateSubscriptionPlan creates a new subscription plan with periods
func (s *BillingService) CreateSubscriptionPlan(ctx context.Context, planData *helper.SubscriptionPlanCreate) error {
	// Check if plan already exists
	existingPlan, err := s.planRepo.GetByTier(ctx, planData.Tier)
	if err != nil {
		return err
	}
	if existingPlan != nil {
		return fmt.Errorf("plan with tier %s already exists", planData.Tier)
	}

	// Create plan
	isActive := true
	if planData.IsActive != nil {
		isActive = *planData.IsActive
	}

	plan := &models.SubscriptionPlan{
		Tier:     planData.Tier,
		Name:     planData.Name,
		Category: planData.Category,
		IsActive: isActive,
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return err
	}

	// Create periods
	for _, periodData := range planData.Periods {
		period := &models.SubscriptionPlanPeriod{
			PlanTier:          planData.Tier,
			Period:            periodData.Period,
			Credits:           periodData.Credits,
			RatePerCredit:     periodData.RatePerCredit,
			Price:             periodData.Price,
			SavingsAmount:     periodData.SavingsAmount,
			SavingsPercentage: periodData.SavingsPercentage,
		}

		if err := s.periodRepo.Create(ctx, period); err != nil {
			return err
		}
	}

	return nil
}

// UpdateSubscriptionPlan updates a subscription plan
func (s *BillingService) UpdateSubscriptionPlan(ctx context.Context, tier string, updateData *helper.SubscriptionPlanUpdate) error {
	plan, err := s.planRepo.GetByTier(ctx, tier)
	if err != nil {
		return err
	}
	if plan == nil {
		return fmt.Errorf("plan with tier %s not found", tier)
	}

	if updateData.Name != nil {
		plan.Name = *updateData.Name
	}
	if updateData.Category != nil {
		plan.Category = *updateData.Category
	}
	if updateData.IsActive != nil {
		plan.IsActive = *updateData.IsActive
	}

	return s.planRepo.Update(ctx, plan)
}

// DeleteSubscriptionPlan deletes a subscription plan
func (s *BillingService) DeleteSubscriptionPlan(ctx context.Context, tier string) error {
	plan, err := s.planRepo.GetByTier(ctx, tier)
	if err != nil {
		return err
	}
	if plan == nil {
		return fmt.Errorf("plan with tier %s not found", tier)
	}

	return s.planRepo.Delete(ctx, tier)
}

// CreateSubscriptionPlanPeriod creates or updates a subscription plan period
func (s *BillingService) CreateSubscriptionPlanPeriod(ctx context.Context, tier string, periodData *helper.SubscriptionPeriodCreate) error {
	plan, err := s.planRepo.GetByTier(ctx, tier)
	if err != nil {
		return err
	}
	if plan == nil {
		return fmt.Errorf("plan with tier %s not found", tier)
	}

	// Check if period exists
	existingPeriod, err := s.periodRepo.GetByPlanAndPeriod(ctx, tier, periodData.Period)
	if err != nil {
		return err
	}

	if existingPeriod != nil {
		// Update existing period
		existingPeriod.Credits = periodData.Credits
		existingPeriod.RatePerCredit = periodData.RatePerCredit
		existingPeriod.Price = periodData.Price
		existingPeriod.SavingsAmount = periodData.SavingsAmount
		existingPeriod.SavingsPercentage = periodData.SavingsPercentage
		return s.periodRepo.Update(ctx, existingPeriod)
	}

	// Create new period
	period := &models.SubscriptionPlanPeriod{
		PlanTier:          tier,
		Period:            periodData.Period,
		Credits:           periodData.Credits,
		RatePerCredit:     periodData.RatePerCredit,
		Price:             periodData.Price,
		SavingsAmount:     periodData.SavingsAmount,
		SavingsPercentage: periodData.SavingsPercentage,
	}

	return s.periodRepo.Create(ctx, period)
}

// DeleteSubscriptionPlanPeriod deletes a subscription plan period
func (s *BillingService) DeleteSubscriptionPlanPeriod(ctx context.Context, tier, period string) error {
	periodObj, err := s.periodRepo.GetByPlanAndPeriod(ctx, tier, period)
	if err != nil {
		return err
	}
	if periodObj == nil {
		return fmt.Errorf("period %s not found for plan %s", period, tier)
	}

	return s.periodRepo.Delete(ctx, periodObj.ID)
}

// Admin CRUD methods for Addon Packages

// CreateAddonPackage creates a new addon package
func (s *BillingService) CreateAddonPackage(ctx context.Context, packageData *helper.AddonPackageCreate) error {
	// Check if package already exists
	existingPackage, err := s.addonRepo.GetByID(ctx, packageData.ID)
	if err != nil {
		return err
	}
	if existingPackage != nil {
		return fmt.Errorf("package with id %s already exists", packageData.ID)
	}

	isActive := true
	if packageData.IsActive != nil {
		isActive = *packageData.IsActive
	}

	pkg := &models.AddonPackage{
		ID:            packageData.ID,
		Name:          packageData.Name,
		Credits:       packageData.Credits,
		RatePerCredit: packageData.RatePerCredit,
		Price:         packageData.Price,
		IsActive:      isActive,
	}

	return s.addonRepo.Create(ctx, pkg)
}

// UpdateAddonPackage updates an addon package
func (s *BillingService) UpdateAddonPackage(ctx context.Context, packageID string, updateData *helper.AddonPackageUpdate) error {
	pkg, err := s.addonRepo.GetByID(ctx, packageID)
	if err != nil {
		return err
	}
	if pkg == nil {
		return fmt.Errorf("package with id %s not found", packageID)
	}

	if updateData.Name != nil {
		pkg.Name = *updateData.Name
	}
	if updateData.Credits != nil {
		pkg.Credits = *updateData.Credits
	}
	if updateData.RatePerCredit != nil {
		pkg.RatePerCredit = *updateData.RatePerCredit
	}
	if updateData.Price != nil {
		pkg.Price = *updateData.Price
	}
	if updateData.IsActive != nil {
		pkg.IsActive = *updateData.IsActive
	}

	return s.addonRepo.Update(ctx, pkg)
}

// DeleteAddonPackage deletes an addon package
func (s *BillingService) DeleteAddonPackage(ctx context.Context, packageID string) error {
	pkg, err := s.addonRepo.GetByID(ctx, packageID)
	if err != nil {
		return err
	}
	if pkg == nil {
		return fmt.Errorf("package with id %s not found", packageID)
	}

	return s.addonRepo.Delete(ctx, packageID)
}

