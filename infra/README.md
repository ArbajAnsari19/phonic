# Infrastructure

This directory contains infrastructure definitions and deployment configurations.

## Directories
- `docker/` - Dockerfile definitions for each service
- `k8s/` - Kubernetes manifests for production deployment

## Local Development
- Use `docker-compose.yml` (in project root) for local development
- Includes: Redis, PostgreSQL, MinIO, Moshi STT/TTS servers

## Production
- Kubernetes manifests for scalable deployment
- Helm charts for configuration management
