-- Create user_profiles table
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id TEXT UNIQUE NOT NULL,
    job_title VARCHAR(255),
    bio TEXT,
    timezone VARCHAR(100),
    avatar_url TEXT,
    notifications JSONB DEFAULT '{}',
    role VARCHAR(50) DEFAULT 'Member',
    credits INTEGER NOT NULL DEFAULT 0,
    subscription_plan VARCHAR(50) DEFAULT 'free',
    subscription_period VARCHAR(20) DEFAULT 'monthly',
    subscription_status VARCHAR(50) DEFAULT 'active',
    subscription_started_at TIMESTAMPTZ,
    subscription_ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users(uuid) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);

