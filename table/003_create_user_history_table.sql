-- Create user_history table
CREATE TABLE IF NOT EXISTS user_history (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    ip VARCHAR(45),
    device TEXT,
    continent VARCHAR(50),
    continent_code VARCHAR(2),
    country VARCHAR(100),
    country_code VARCHAR(2),
    region VARCHAR(10),
    region_name VARCHAR(100),
    city VARCHAR(100),
    district VARCHAR(100),
    zip VARCHAR(20),
    lat NUMERIC(10, 7),
    lon NUMERIC(10, 7),
    timezone VARCHAR(100),
    currency VARCHAR(10),
    isp VARCHAR(255),
    org VARCHAR(255),
    asname VARCHAR(255),
    reverse VARCHAR(255),
    proxy BOOLEAN DEFAULT false,
    hosting BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_history_user_id FOREIGN KEY (user_id) REFERENCES users(uuid) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_user_history_user_id ON user_history(user_id);
CREATE INDEX IF NOT EXISTS idx_user_history_event_type ON user_history(event_type);
CREATE INDEX IF NOT EXISTS idx_user_history_created_at ON user_history(created_at);

