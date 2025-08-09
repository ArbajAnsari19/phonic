# Health Checks

This document explains the health check system for the Phonic AI Calling Agent.

## Overview

The health check system provides comprehensive monitoring of service health and dependencies. It's designed to work with Kubernetes, load balancers, and monitoring systems.

## Health Check Types

### 1. Liveness Checks
- **Purpose**: Determine if the service is running
- **Endpoint**: `/health/live`
- **Kubernetes**: Used for liveness probes
- **Behavior**: Simple check that returns 200 if service is alive

### 2. Readiness Checks
- **Purpose**: Determine if the service is ready to handle requests
- **Endpoint**: `/health/ready`
- **Kubernetes**: Used for readiness probes
- **Behavior**: Returns 200 only if all dependencies are healthy

### 3. Full Health Checks
- **Purpose**: Comprehensive health status with details
- **Endpoint**: `/health`
- **Monitoring**: Used by Prometheus and monitoring dashboards
- **Behavior**: Returns detailed JSON with all check results

## Health Check Components

### Database Checker
- **Name**: `database`
- **Checks**: PostgreSQL connectivity and connection pool status
- **Metadata**: Connection counts, pool statistics
- **Timeout**: 5 seconds

### Redis Checker
- **Name**: `redis`
- **Checks**: Redis connectivity and basic operations
- **Metadata**: Ping response, server info availability
- **Timeout**: 5 seconds

### Moshi STT Checker
- **Name**: `moshi_stt`
- **Checks**: Moshi STT server HTTP health endpoint
- **Metadata**: Service URL, response status
- **Timeout**: 5 seconds

### Moshi TTS Checker
- **Name**: `moshi_tts`
- **Checks**: Moshi TTS server HTTP health endpoint
- **Metadata**: Service URL, response status
- **Timeout**: 5 seconds

### Custom Checkers
- **Purpose**: Service-specific health checks
- **Examples**: System resources, file permissions, external APIs
- **Implementation**: Flexible function-based checks

## Health Check Responses

### Successful Response (200 OK)
```json
{
  "status": "healthy",
  "timestamp": "2025-08-09T18:54:34+05:30",
  "service": "Phonic AI Calling Agent",
  "version": "0.1.0-dev",
  "uptime": "5m30s",
  "checks": [
    {
      "name": "database",
      "status": "healthy",
      "message": "Database connection healthy",
      "duration": "340µs",
      "metadata": {
        "open_connections": "1",
        "in_use": "0",
        "idle": "1"
      },
      "timestamp": "2025-08-09T18:54:34+05:30"
    }
  ]
}
```

### Failed Response (503 Service Unavailable)
```json
{
  "status": "unhealthy",
  "timestamp": "2025-08-09T18:54:34+05:30",
  "service": "Phonic AI Calling Agent",
  "version": "0.1.0-dev",
  "uptime": "5m30s",
  "checks": [
    {
      "name": "moshi_stt",
      "status": "unhealthy",
      "message": "Moshi service unreachable: connection refused",
      "duration": "5s",
      "timestamp": "2025-08-09T18:54:34+05:30"
    }
  ]
}
```

## Kubernetes Configuration

### Deployment Example
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: phonic-gateway
spec:
  template:
    spec:
      containers:
      - name: gateway
        image: phonic/gateway:latest
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
```

## Testing Health Checks

### Manual Testing
```bash
# Test all health checks
make health-test

# Test specific endpoints
curl http://localhost:8080/health
curl http://localhost:8080/health/ready
curl http://localhost:8080/health/live
```

### Load Balancer Configuration

#### Nginx
```nginx
upstream phonic_backend {
    server phonic-1:8080;
    server phonic-2:8080;
}

server {
    location /health/ready {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
    
    location / {
        proxy_pass http://phonic_backend;
        proxy_next_upstream error timeout http_502 http_503 http_504;
    }
}
```

#### HAProxy
```
backend phonic_servers
    balance roundrobin
    option httpchk GET /health/ready
    server phonic-1 phonic-1:8080 check
    server phonic-2 phonic-2:8080 check
```

## Monitoring Integration

### Prometheus Metrics
Health check results can be exposed as Prometheus metrics:

```
# Health check status (1 = healthy, 0 = unhealthy)
phonic_health_check_status{service="phonic-gateway", check="database"} 1

# Health check duration in seconds
phonic_health_check_duration_seconds{service="phonic-gateway", check="database"} 0.0003405
```

### Grafana Dashboard
Create dashboards to visualize:
- Overall service health status
- Individual dependency health
- Health check response times
- Health check failure rates

## Graceful Shutdown

The health check system integrates with graceful shutdown:

1. **Signal Reception**: Service receives SIGTERM/SIGINT
2. **Readiness**: Readiness checks start failing immediately
3. **Graceful Period**: Service completes in-flight requests
4. **Cleanup**: All shutdown hooks execute
5. **Termination**: Service exits cleanly

### Shutdown Hooks
```go
shutdownManager := shutdown.NewManager(30*time.Second, logger)

// Add cleanup hooks
shutdownManager.AddHook(func(ctx context.Context) error {
    return database.Close()
})

shutdownManager.AddHook(func(ctx context.Context) error {
    return httpServer.Shutdown(ctx)
})

// Wait for shutdown signal
shutdownManager.WaitForShutdown()
```

## Best Practices

1. **Timeout Configuration**: Set appropriate timeouts for each check
2. **Dependency Ordering**: Check critical dependencies first
3. **Graceful Degradation**: Service can still function with some dependencies down
4. **Monitoring**: Alert on health check failures
5. **Testing**: Regularly test health checks in different failure scenarios

## Troubleshooting

### Common Issues

1. **Slow Health Checks**: Increase timeouts or optimize check logic
2. **False Negatives**: Dependencies temporarily unavailable
3. **Resource Exhaustion**: Health checks consuming too many resources
4. **Network Issues**: DNS resolution or connectivity problems

### Debug Commands
```bash
# Test health checks with verbose output
go run cmd/health-test/main.go dev

# Check specific dependency
curl -v http://localhost:8080/health

# Monitor health check logs
docker logs phonic-gateway | grep "health check"
```
