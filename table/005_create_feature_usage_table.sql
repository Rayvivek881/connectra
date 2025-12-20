-- Create feature_usage table
CREATE TABLE IF NOT EXISTS feature_usage (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    feature VARCHAR(50) NOT NULL,
    used INTEGER NOT NULL DEFAULT 0,
    "limit" INTEGER NOT NULL DEFAULT 0,
    period_start TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    period_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_feature_usage_user_id FOREIGN KEY (user_id) REFERENCES users(uuid) ON DELETE CASCADE,
    CONSTRAINT unique_user_feature UNIQUE (user_id, feature)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_feature_usage_user_id ON feature_usage(user_id);
CREATE INDEX IF NOT EXISTS idx_feature_usage_feature ON feature_usage(feature);
CREATE INDEX IF NOT EXISTS idx_feature_usage_user_feature ON feature_usage(user_id, feature);
CREATE INDEX IF NOT EXISTS idx_feature_usage_period_start ON feature_usage(period_start);

