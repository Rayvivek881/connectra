-- Create subscription_plans table
CREATE TABLE IF NOT EXISTS subscription_plans (
    tier VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL, -- STARTER, PROFESSIONAL, BUSINESS, ENTERPRISE
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Create subscription_plan_periods table
CREATE TABLE IF NOT EXISTS subscription_plan_periods (
    id SERIAL PRIMARY KEY,
    plan_tier VARCHAR(50) NOT NULL REFERENCES subscription_plans(tier) ON DELETE CASCADE,
    period VARCHAR(20) NOT NULL, -- monthly, quarterly, yearly
    credits INTEGER NOT NULL,
    rate_per_credit NUMERIC(10, 6) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    savings_amount NUMERIC(10, 2),
    savings_percentage INTEGER, -- Percentage as integer (e.g., 10 for 10%)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    UNIQUE(plan_tier, period)
);

-- Create addon_packages table
CREATE TABLE IF NOT EXISTS addon_packages (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    credits INTEGER NOT NULL,
    rate_per_credit NUMERIC(10, 6) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_subscription_plans_tier ON subscription_plans(tier);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_category ON subscription_plans(category);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_is_active ON subscription_plans(is_active);

CREATE INDEX IF NOT EXISTS idx_subscription_plan_periods_plan_tier ON subscription_plan_periods(plan_tier);
CREATE INDEX IF NOT EXISTS idx_subscription_plan_periods_period ON subscription_plan_periods(period);

CREATE INDEX IF NOT EXISTS idx_addon_packages_id ON addon_packages(id);
CREATE INDEX IF NOT EXISTS idx_addon_packages_is_active ON addon_packages(is_active);

-- Seed subscription plans
INSERT INTO subscription_plans (tier, name, category, is_active) VALUES
    ('5k', '5k Credits Tier', 'STARTER', true),
    ('25k', '25k Credits Tier', 'STARTER', true),
    ('100k', '100k Credits Tier', 'PROFESSIONAL', true),
    ('500k', '500k Credits Tier', 'PROFESSIONAL', true),
    ('1M', '1M Credits Tier', 'BUSINESS', true),
    ('5M', '5M Credits Tier', 'BUSINESS', true),
    ('10M', '10M Credits Tier', 'ENTERPRISE', true)
ON CONFLICT (tier) DO NOTHING;

-- Seed subscription plan periods
INSERT INTO subscription_plan_periods (plan_tier, period, credits, rate_per_credit, price, savings_amount, savings_percentage) VALUES
    -- 5k tier
    ('5k', 'monthly', 5000, 0.002, 10.0, NULL, NULL),
    ('5k', 'quarterly', 15000, 0.0018, 27.0, 3.0, 10),
    ('5k', 'yearly', 60000, 0.0016, 96.0, 24.0, 20),
    
    -- 25k tier
    ('25k', 'monthly', 25000, 0.0012, 30.0, NULL, NULL),
    ('25k', 'quarterly', 75000, 0.00108, 81.0, 9.0, 10),
    ('25k', 'yearly', 300000, 0.00096, 288.0, 72.0, 20),
    
    -- 100k tier
    ('100k', 'monthly', 100000, 0.00099, 99.0, NULL, NULL),
    ('100k', 'quarterly', 300000, 0.000891, 267.0, 30.0, 10),
    ('100k', 'yearly', 1200000, 0.000792, 950.0, 238.0, 20),
    
    -- 500k tier
    ('500k', 'monthly', 500000, 0.000398, 199.0, NULL, NULL),
    ('500k', 'quarterly', 1500000, 0.0003582, 537.0, 60.0, 10),
    ('500k', 'yearly', 6000000, 0.0003184, 1910.0, 478.0, 20),
    
    -- 1M tier
    ('1M', 'monthly', 1000000, 0.000299, 299.0, NULL, NULL),
    ('1M', 'quarterly', 3000000, 0.0002691, 807.0, 90.0, 10),
    ('1M', 'yearly', 12000000, 0.0002392, 2870.0, 718.0, 20),
    
    -- 5M tier
    ('5M', 'monthly', 5000000, 0.0001998, 999.0, NULL, NULL),
    ('5M', 'quarterly', 15000000, 0.00017982, 2697.0, 300.0, 10),
    ('5M', 'yearly', 60000000, 0.00015984, 9590.0, 2398.0, 20),
    
    -- 10M tier
    ('10M', 'monthly', 10000000, 0.0001599, 1599.0, NULL, NULL),
    ('10M', 'quarterly', 30000000, 0.00014391, 4317.0, 480.0, 10),
    ('10M', 'yearly', 120000000, 0.00012792, 15350.0, 3838.0, 20)
ON CONFLICT (plan_tier, period) DO NOTHING;

-- Seed addon packages
INSERT INTO addon_packages (id, name, credits, rate_per_credit, price, is_active) VALUES
    ('small', 'Small', 5000, 0.002, 10.0, true),
    ('basic', 'Basic', 25000, 0.0012, 30.0, true),
    ('standard', 'Standard', 100000, 0.00099, 99.0, true),
    ('plus', 'Plus', 500000, 0.000398, 199.0, true),
    ('pro', 'Pro', 1000000, 0.000299, 299.0, true),
    ('advanced', 'Advanced', 5000000, 0.0001998, 999.0, true),
    ('premium', 'Premium', 10000000, 0.0001599, 1599.0, true)
ON CONFLICT (id) DO NOTHING;

