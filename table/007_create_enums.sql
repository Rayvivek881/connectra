-- Create PostgreSQL ENUM types for better type safety
-- Note: These are optional as we're using VARCHAR with application-level validation

-- User history event types
DO $$ BEGIN
    CREATE TYPE user_history_event_type AS ENUM ('registration', 'login');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Activity service types
DO $$ BEGIN
    CREATE TYPE activity_service_type AS ENUM ('linkedin', 'email');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Activity action types
DO $$ BEGIN
    CREATE TYPE activity_action_type AS ENUM ('search', 'export');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Activity status types
DO $$ BEGIN
    CREATE TYPE activity_status AS ENUM ('success', 'failed', 'partial');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Feature types
DO $$ BEGIN
    CREATE TYPE feature_type AS ENUM (
        'AI_CHAT',
        'BULK_EXPORT',
        'API_KEYS',
        'TEAM_MANAGEMENT',
        'EMAIL_FINDER',
        'VERIFIER',
        'LINKEDIN',
        'DATA_SEARCH',
        'ADVANCED_FILTERS',
        'AI_SUMMARIES',
        'SAVE_SEARCHES',
        'BULK_VERIFICATION'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

