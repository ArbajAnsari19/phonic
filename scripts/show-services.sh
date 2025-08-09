#!/bin/bash

# Show URLs for all Phonic services

echo "🎵 Phonic AI Calling Agent - Service URLs"
echo "========================================"
echo

echo "📊 Monitoring & Observability:"
echo "   Grafana Dashboard:     http://localhost:3000 (admin/admin)"
echo "   Prometheus Metrics:    http://localhost:9090"
echo

echo "🗄️ Data Storage:"
echo "   MinIO Console:         http://localhost:9001 (phonic/phonic_dev_password)"
echo "   PostgreSQL:            localhost:5432 (phonic/phonic_dev_password)"
echo "   Redis:                 localhost:6379"
echo

echo "🤖 AI Services (Placeholders):"
echo "   Moshi STT Server:      http://localhost:8001"
echo "   Moshi TTS Server:      http://localhost:8002"
echo

echo "🚀 Phonic Services (Coming Soon):"
echo "   Gateway Service:       http://localhost:8080 (Step 45+)"
echo "   Session Service:       http://localhost:8083 (Step 53+)"
echo "   Orchestrator Service:  http://localhost:8084 (Step 59+)"
echo

echo "🛠️ Development Tools (if using docker-compose.dev.yml):"
echo "   Adminer (DB Admin):    http://localhost:8080"
echo "   Redis Commander:       http://localhost:8081"
echo "   File Browser:          http://localhost:8082"
echo "   MailHog:               http://localhost:8025"
echo

echo "📋 Quick Commands:"
echo "   make compose-logs      # View all service logs"
echo "   make compose-down      # Stop all services"
echo "   make compose-restart   # Restart all services"
echo "   docker ps              # Check service status"
echo
