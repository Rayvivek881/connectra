-- Create user_activities table
CREATE TABLE IF NOT EXISTS user_activities (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    service_type VARCHAR(50) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    request_params JSONB,
    result_count INTEGER NOT NULL DEFAULT 0,
    result_summary JSONB,
    error_message TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_activities_user_id FOREIGN KEY (user_id) REFERENCES users(uuid) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_user_activities_user_id ON user_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activities_service_type ON user_activities(service_type);
CREATE INDEX IF NOT EXISTS idx_user_activities_action_type ON user_activities(action_type);
CREATE INDEX IF NOT EXISTS idx_user_activities_status ON user_activities(status);
CREATE INDEX IF NOT EXISTS idx_user_activities_created_at ON user_activities(created_at);
CREATE INDEX IF NOT EXISTS idx_user_activities_user_service_action_created ON user_activities(user_id, service_type, action_type, created_at);

