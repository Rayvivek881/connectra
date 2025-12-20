-- Update token_blacklist table to add user_id and update timestamps
ALTER TABLE token_blacklist 
ADD COLUMN IF NOT EXISTS user_id TEXT,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Update expires_at to use TIMESTAMPTZ if not already
ALTER TABLE token_blacklist 
ALTER COLUMN expires_at TYPE TIMESTAMPTZ USING expires_at::TIMESTAMPTZ;

-- Add foreign key constraint if user_id is set
-- Note: This is optional since user_id can be NULL for old tokens

