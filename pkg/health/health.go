// Package health provides health check functionality for Phonic AI Calling Agent
package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/ArbajAnsari19/phonic/pkg/logger"
)

// Status represents the health status of a component
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusUnknown   Status = "unknown"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Name      string            `json:"name"`
	Status    Status            `json:"status"`
	Message   string            `json:"message,omitempty"`
	Duration  time.Duration     `json:"duration"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status    Status        `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Service   string        `json:"service"`
	Version   string        `json:"version"`
	Uptime    time.Duration `json:"uptime"`
	Checks    []CheckResult `json:"checks"`
}

// Checker interface for health check components
type Checker interface {
	Check(ctx context.Context) CheckResult
}

// Manager manages health checks for the service
type Manager struct {
	serviceName string
	version     string
	startTime   time.Time
	checkers    map[string]Checker
	mu          sync.RWMutex
	logger      *logger.Logger
}

// NewManager creates a new health check manager
func NewManager(serviceName, version string, log *logger.Logger) *Manager {
	return &Manager{
		serviceName: serviceName,
		version:     version,
		startTime:   time.Now(),
		checkers:    make(map[string]Checker),
		logger:      log,
	}
}

// AddChecker adds a health checker
func (m *Manager) AddChecker(name string, checker Checker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkers[name] = checker
}

// RemoveChecker removes a health checker
func (m *Manager) RemoveChecker(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.checkers, name)
}

// CheckHealth performs all health checks
func (m *Manager) CheckHealth(ctx context.Context) HealthResponse {
	start := time.Now()
	
	m.mu.RLock()
	checkers := make(map[string]Checker, len(m.checkers))
	for name, checker := range m.checkers {
		checkers[name] = checker
	}
	m.mu.RUnlock()

	var checks []CheckResult
	var wg sync.WaitGroup
	resultCh := make(chan CheckResult, len(checkers))

	// Run all checks concurrently
	for name, checker := range checkers {
		wg.Add(1)
		go func(n string, c Checker) {
			defer wg.Done()
			result := c.Check(ctx)
			result.Name = n
			resultCh <- result
		}(name, checker)
	}

	// Wait for all checks to complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	for result := range resultCh {
		checks = append(checks, result)
	}

	// Determine overall status
	overallStatus := StatusHealthy
	for _, check := range checks {
		if check.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
			break
		}
	}

	duration := time.Since(start)
	m.logger.Info("Health check completed",
		zap.String("overall_status", string(overallStatus)),
		zap.Duration("duration", duration),
		zap.Int("checks_count", len(checks)),
	)

	return HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Service:   m.serviceName,
		Version:   m.version,
		Uptime:    time.Since(m.startTime),
		Checks:    checks,
	}
}

// HTTPHandler returns an HTTP handler for health checks
func (m *Manager) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		health := m.CheckHealth(ctx)

		w.Header().Set("Content-Type", "application/json")
		
		// Set HTTP status code based on health
		if health.Status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err := json.NewEncoder(w).Encode(health); err != nil {
			m.logger.Error("Failed to encode health response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// ReadinessHandler returns a simple readiness check handler
func (m *Manager) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		health := m.CheckHealth(ctx)

		if health.Status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("not ready"))
		}
	}
}

// LivenessHandler returns a simple liveness check handler
func (m *Manager) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Liveness check is simpler - just check if service is running
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("alive"))
	}
}

// DatabaseChecker checks database connectivity
type DatabaseChecker struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewDatabaseChecker creates a new database checker
func NewDatabaseChecker(db *sql.DB, log *logger.Logger) *DatabaseChecker {
	return &DatabaseChecker{
		db:     db,
		logger: log,
	}
}

// Check performs the database health check
func (c *DatabaseChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()
	
	// Simple ping to check connectivity
	err := c.db.PingContext(ctx)
	duration := time.Since(start)
	
	if err != nil {
		c.logger.Error("Database health check failed", zap.Error(err), zap.Duration("duration", duration))
		return CheckResult{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Database ping failed: %v", err),
			Duration:  duration,
			Timestamp: time.Now(),
		}
	}

	// Check database stats
	stats := c.db.Stats()
	metadata := map[string]string{
		"open_connections": fmt.Sprintf("%d", stats.OpenConnections),
		"in_use":          fmt.Sprintf("%d", stats.InUse),
		"idle":            fmt.Sprintf("%d", stats.Idle),
	}

	return CheckResult{
		Status:    StatusHealthy,
		Message:   "Database connection healthy",
		Duration:  duration,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}

// RedisChecker checks Redis connectivity
type RedisChecker struct {
	client *redis.Client
	logger *logger.Logger
}

// NewRedisChecker creates a new Redis checker
func NewRedisChecker(client *redis.Client, log *logger.Logger) *RedisChecker {
	return &RedisChecker{
		client: client,
		logger: log,
	}
}

// Check performs the Redis health check
func (c *RedisChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()
	
	// Simple ping to check connectivity
	pong, err := c.client.Ping(ctx).Result()
	duration := time.Since(start)
	
	if err != nil {
		c.logger.Error("Redis health check failed", zap.Error(err), zap.Duration("duration", duration))
		return CheckResult{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Redis ping failed: %v", err),
			Duration:  duration,
			Timestamp: time.Now(),
		}
	}

	if pong != "PONG" {
		return CheckResult{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Redis ping returned unexpected response: %s", pong),
			Duration:  duration,
			Timestamp: time.Now(),
		}
	}

	// Get Redis info
	_, err = c.client.Info(ctx).Result()
	metadata := map[string]string{
		"ping_response": pong,
	}
	if err == nil {
		metadata["info_available"] = "true"
	}

	return CheckResult{
		Status:    StatusHealthy,
		Message:   "Redis connection healthy",
		Duration:  duration,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}

// MoshiChecker checks Moshi STT/TTS server connectivity
type MoshiChecker struct {
	serviceURL string
	timeout    time.Duration
	logger     *logger.Logger
}

// NewMoshiChecker creates a new Moshi service checker
func NewMoshiChecker(serviceURL string, timeout time.Duration, log *logger.Logger) *MoshiChecker {
	return &MoshiChecker{
		serviceURL: serviceURL,
		timeout:    timeout,
		logger:     log,
	}
}

// Check performs the Moshi service health check
func (c *MoshiChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: c.timeout,
	}
	
	// Simple HTTP GET to check if service is responding
	// Note: This is a placeholder - actual Moshi health endpoint may be different
	healthURL := fmt.Sprintf("http://%s/health", c.serviceURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		duration := time.Since(start)
		return CheckResult{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Failed to create request: %v", err),
			Duration:  duration,
			Timestamp: time.Now(),
		}
	}
	
	resp, err := client.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		c.logger.Error("Moshi health check failed", zap.Error(err), zap.Duration("duration", duration))
		return CheckResult{
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Moshi service unreachable: %v", err),
			Duration:  duration,
			Timestamp: time.Now(),
		}
	}
	defer resp.Body.Close()

	metadata := map[string]string{
		"service_url":   c.serviceURL,
		"status_code":   fmt.Sprintf("%d", resp.StatusCode),
		"content_type":  resp.Header.Get("Content-Type"),
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return CheckResult{
			Status:    StatusHealthy,
			Message:   "Moshi service responding",
			Duration:  duration,
			Metadata:  metadata,
			Timestamp: time.Now(),
		}
	}

	return CheckResult{
		Status:    StatusUnhealthy,
		Message:   fmt.Sprintf("Moshi service returned status %d", resp.StatusCode),
		Duration:  duration,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}

// CustomChecker allows for custom health checks
type CustomChecker struct {
	name      string
	checkFunc func(ctx context.Context) (bool, string, map[string]string)
	logger    *logger.Logger
}

// NewCustomChecker creates a new custom checker
func NewCustomChecker(name string, checkFunc func(ctx context.Context) (bool, string, map[string]string), log *logger.Logger) *CustomChecker {
	return &CustomChecker{
		name:      name,
		checkFunc: checkFunc,
		logger:    log,
	}
}

// Check performs the custom health check
func (c *CustomChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()
	
	healthy, message, metadata := c.checkFunc(ctx)
	duration := time.Since(start)
	
	status := StatusHealthy
	if !healthy {
		status = StatusUnhealthy
	}

	return CheckResult{
		Status:    status,
		Message:   message,
		Duration:  duration,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}
