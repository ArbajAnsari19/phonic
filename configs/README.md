# Configurations

Environment-specific configuration files for the Phonic AI Calling Agent.

## Environments
- `dev/` - Development environment settings
- `staging/` - Staging environment settings  
- `prod/` - Production environment settings

## Configuration Files
- `app.yaml` - Main application settings
- `database.yaml` - Database connection settings
- `redis.yaml` - Redis configuration
- `moshi.yaml` - Moshi STT/TTS server settings
- `logging.yaml` - Logging configuration

## Usage
Configurations are loaded based on `PHONIC_ENV` environment variable.
