-- Phonic AI Calling Agent - Development Database Setup
-- Extended schema with development-specific features

-- Include base schema
\i /docker-entrypoint-initdb.d/init.sql

-- Development-specific configurations
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = 'on';
ALTER SYSTEM SET log_min_duration_statement = 0;

-- Create development admin user
CREATE USER phonic_admin WITH PASSWORD 'admin_dev_password';
GRANT ALL PRIVILEGES ON DATABASE phonic_dev TO phonic_admin;
ALTER USER phonic_admin CREATEDB;
