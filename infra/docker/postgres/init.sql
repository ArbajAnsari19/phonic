-- Phonic AI Calling Agent - PostgreSQL Initialization
-- Production database schema setup

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS analytics;

-- Set default search path
ALTER DATABASE phonic SET search_path TO public, analytics;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE
);

-- Create call sessions table
CREATE TABLE IF NOT EXISTS call_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    session_token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create call events table (for analytics)
CREATE TABLE IF NOT EXISTS analytics.call_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES call_sessions(id),
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB DEFAULT '{}',
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    service_name VARCHAR(100),
    trace_id VARCHAR(100)
);

-- Create audio files table
CREATE TABLE IF NOT EXISTS audio_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES call_sessions(id),
    file_type VARCHAR(20) NOT NULL, -- 'input', 'output', 'processed'
    file_path VARCHAR(500) NOT NULL, -- MinIO path
    file_size_bytes BIGINT,
    duration_seconds FLOAT,
    format VARCHAR(20), -- 'wav', 'mp3', 'opus'
    sample_rate INTEGER,
    channels INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_call_sessions_user_id ON call_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_call_sessions_status ON call_sessions(status);
CREATE INDEX IF NOT EXISTS idx_call_sessions_started_at ON call_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_call_events_session_id ON analytics.call_events(session_id);
CREATE INDEX IF NOT EXISTS idx_call_events_timestamp ON analytics.call_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_call_events_event_type ON analytics.call_events(event_type);
CREATE INDEX IF NOT EXISTS idx_audio_files_session_id ON audio_files(session_id);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_call_sessions_updated_at BEFORE UPDATE ON call_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO phonic;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA analytics TO phonic;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO phonic;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA analytics TO phonic;

-- Create a test user for development
INSERT INTO users (email, name) VALUES 
    ('test@phonic.ai', 'Test User')
ON CONFLICT (email) DO NOTHING;