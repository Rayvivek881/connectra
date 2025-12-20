package models

import (
	"context"
	"vivek-ray/connections"

	"github.com/rs/zerolog/log"
)

// SubscriptionPlanRepository provides data access methods for subscription plans
type SubscriptionPlanRepository struct{}

func NewSubscriptionPlanRepository() *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{}
}

// GetByTier retrieves a subscription plan by tier
func (r *SubscriptionPlanRepository) GetByTier(ctx context.Context, tier string) (*SubscriptionPlan, error) {
	plan := new(SubscriptionPlan)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(plan).
		Where("tier = ?", tier).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding subscription plan by tier")
		return nil, err
	}
	return plan, nil
}

// ListAll retrieves all subscription plans
func (r *SubscriptionPlanRepository) ListAll(ctx context.Context, includeInactive bool) ([]SubscriptionPlan, error) {
	var plans []SubscriptionPlan
	query := connections.PgDBConnection.Client.NewSelect().Model(&plans)
	
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	
	query = query.Order("tier ASC")
	
	err := query.Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing subscription plans")
		return nil, err
	}
	return plans, nil
}

// ListAllWithPeriods retrieves all subscription plans with their periods (optimized with JOIN)
func (r *SubscriptionPlanRepository) ListAllWithPeriods(ctx context.Context, includeInactive bool) ([]SubscriptionPlan, error) {
	var plans []SubscriptionPlan
	query := connections.PgDBConnection.Client.NewSelect().
		Model(&plans).
		Relation("Periods")
	
	if !includeInactive {
		query = query.Where("sp.is_active = ?", true)
	}
	
	query = query.Order("tier ASC")
	
	err := query.Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing subscription plans with periods")
		return nil, err
	}
	return plans, nil
}

// Create creates a new subscription plan
func (r *SubscriptionPlanRepository) Create(ctx context.Context, plan *SubscriptionPlan) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(plan).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating subscription plan")
		return err
	}
	return nil
}

// Update updates a subscription plan
func (r *SubscriptionPlanRepository) Update(ctx context.Context, plan *SubscriptionPlan) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model(plan).
		Where("tier = ?", plan.Tier).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error updating subscription plan")
		return err
	}
	return nil
}

// Delete deletes a subscription plan (cascade will delete periods)
func (r *SubscriptionPlanRepository) Delete(ctx context.Context, tier string) error {
	_, err := connections.PgDBConnection.Client.NewDelete().
		Model((*SubscriptionPlan)(nil)).
		Where("tier = ?", tier).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting subscription plan")
		return err
	}
	return nil
}

// SubscriptionPlanPeriodRepository provides data access methods for subscription plan periods
type SubscriptionPlanPeriodRepository struct{}

func NewSubscriptionPlanPeriodRepository() *SubscriptionPlanPeriodRepository {
	return &SubscriptionPlanPeriodRepository{}
}

// GetByPlanAndPeriod retrieves a period by plan tier and period
func (r *SubscriptionPlanPeriodRepository) GetByPlanAndPeriod(ctx context.Context, planTier, period string) (*SubscriptionPlanPeriod, error) {
	periodObj := new(SubscriptionPlanPeriod)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(periodObj).
		Where("plan_tier = ? AND period = ?", planTier, period).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding subscription plan period")
		return nil, err
	}
	return periodObj, nil
}

// ListByPlan retrieves all periods for a subscription plan
func (r *SubscriptionPlanPeriodRepository) ListByPlan(ctx context.Context, planTier string) ([]SubscriptionPlanPeriod, error) {
	var periods []SubscriptionPlanPeriod
	err := connections.PgDBConnection.Client.NewSelect().
		Model(&periods).
		Where("plan_tier = ?", planTier).
		Order("period ASC").
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing subscription plan periods")
		return nil, err
	}
	return periods, nil
}

// Create creates a new subscription plan period
func (r *SubscriptionPlanPeriodRepository) Create(ctx context.Context, period *SubscriptionPlanPeriod) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(period).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating subscription plan period")
		return err
	}
	return nil
}

// Update updates a subscription plan period
func (r *SubscriptionPlanPeriodRepository) Update(ctx context.Context, period *SubscriptionPlanPeriod) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model(period).
		Where("id = ?", period.ID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error updating subscription plan period")
		return err
	}
	return nil
}

// Delete deletes a subscription plan period
func (r *SubscriptionPlanPeriodRepository) Delete(ctx context.Context, id int) error {
	_, err := connections.PgDBConnection.Client.NewDelete().
		Model((*SubscriptionPlanPeriod)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting subscription plan period")
		return err
	}
	return nil
}

// AddonPackageRepository provides data access methods for addon packages
type AddonPackageRepository struct{}

func NewAddonPackageRepository() *AddonPackageRepository {
	return &AddonPackageRepository{}
}

// GetByID retrieves an addon package by ID
func (r *AddonPackageRepository) GetByID(ctx context.Context, packageID string) (*AddonPackage, error) {
	pkg := new(AddonPackage)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(pkg).
		Where("id = ?", packageID).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding addon package by ID")
		return nil, err
	}
	return pkg, nil
}

// ListAll retrieves all addon packages
func (r *AddonPackageRepository) ListAll(ctx context.Context, includeInactive bool) ([]AddonPackage, error) {
	var packages []AddonPackage
	query := connections.PgDBConnection.Client.NewSelect().Model(&packages)
	
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	
	query = query.Order("price ASC")
	
	err := query.Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing addon packages")
		return nil, err
	}
	return packages, nil
}

// Create creates a new addon package
func (r *AddonPackageRepository) Create(ctx context.Context, pkg *AddonPackage) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(pkg).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating addon package")
		return err
	}
	return nil
}

// Update updates an addon package
func (r *AddonPackageRepository) Update(ctx context.Context, pkg *AddonPackage) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model(pkg).
		Where("id = ?", pkg.ID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error updating addon package")
		return err
	}
	return nil
}

// Delete deletes an addon package
func (r *AddonPackageRepository) Delete(ctx context.Context, packageID string) error {
	_, err := connections.PgDBConnection.Client.NewDelete().
		Model((*AddonPackage)(nil)).
		Where("id = ?", packageID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting addon package")
		return err
	}
	return nil
}

