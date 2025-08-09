# Configuration Management

This document explains how to configure the Phonic AI Calling Agent for different environments.

## Configuration Files

The application uses YAML configuration files located in the `configs/` directory:

- `configs/dev/app.yaml` - Development environment
- `configs/staging/app.yaml` - Staging environment  
- `configs/prod/app.yaml` - Production environment

## Environment Selection

The environment is selected using the `PHONIC_ENV` environment variable:

```bash
export PHONIC_ENV=dev     # Development (default)
export PHONIC_ENV=staging # Staging
export PHONIC_ENV=prod    # Production
```

## Environment Variables

All configuration values can be overridden using environment variables with the `PHONIC_` prefix. Variable names use underscores instead of dots:

### Application Settings
```bash
PHONIC_APP_NAME="Phonic AI Calling Agent"
PHONIC_APP_ENVIRONMENT=dev
PHONIC_APP_DEBUG=true
PHONIC_APP_PORT=8080
```

### Database Configuration
```bash
PHONIC_DATABASE_HOST=localhost
PHONIC_DATABASE_PORT=5432
PHONIC_DATABASE_USERNAME=phonic
PHONIC_DATABASE_PASSWORD=your_password
PHONIC_DATABASE_DATABASE=phonic
PHONIC_DATABASE_SSL_MODE=disable
```

### Redis Configuration
```bash
PHONIC_REDIS_HOST=localhost
PHONIC_REDIS_PORT=6379
PHONIC_REDIS_PASSWORD=your_password
PHONIC_REDIS_DATABASE=0
```

### Moshi STT/TTS Configuration
```bash
PHONIC_MOSHI_STT_HOST=localhost
PHONIC_MOSHI_STT_PORT=8001
PHONIC_MOSHI_TTS_HOST=localhost
PHONIC_MOSHI_TTS_PORT=8002
```

### Security Configuration
```bash
PHONIC_SECURITY_JWT_SECRET=your-secret-key
PHONIC_SECURITY_JWT_EXPIRY_HOURS=24
PHONIC_SECURITY_RATE_LIMIT_REQUESTS_PER_MINUTE=100
```

### Storage Configuration (MinIO/S3)
```bash
PHONIC_STORAGE_ENDPOINT=localhost:9000
PHONIC_STORAGE_ACCESS_KEY=your_access_key
PHONIC_STORAGE_SECRET_KEY=your_secret_key
PHONIC_STORAGE_BUCKET=phonic-audio
PHONIC_STORAGE_REGION=us-east-1
```

## Testing Configuration

Use the config test utility to verify your configuration:

```bash
# Test all environments
make config-test

# Test specific environment
go run cmd/config-test/main.go dev
go run cmd/config-test/main.go staging
go run cmd/config-test/main.go prod
```

## Configuration Validation

The configuration system automatically validates:

- Required fields are present
- Environment values are valid (dev, staging, prod)
- Logging levels are valid (debug, info, warn, error)
- Timeout durations are properly formatted
- URLs and connection strings are formatted correctly

## Development Setup

For local development, the default configuration in `configs/dev/app.yaml` should work with the Docker Compose setup:

```bash
# Start infrastructure services
make compose-up

# Test configuration
make config-test

# All services should be accessible with default settings
```

## Production Deployment

For production deployment:

1. Set `PHONIC_ENV=prod`
2. Provide all required environment variables
3. Ensure database and Redis are configured with SSL
4. Use strong JWT secrets and appropriate rate limits
5. Configure CORS for your frontend domains

## Security Considerations

- Never commit passwords or secrets to version control
- Use environment variables for all sensitive configuration
- Rotate JWT secrets regularly
- Use SSL/TLS for all external connections in production
- Set appropriate rate limits for your use case
