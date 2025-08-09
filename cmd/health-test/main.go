// Health check test utility for Phonic AI Calling Agent
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"

	"github.com/ArbajAnsari19/phonic/pkg/config"
	"github.com/ArbajAnsari19/phonic/pkg/health"
	"github.com/ArbajAnsari19/phonic/pkg/logger"
)

func main() {
	fmt.Println("🎵 Phonic Health Check Test")
	fmt.Println("============================")
	
	// Get environment from command line or default to dev
	environment := "dev"
	if len(os.Args) > 1 {
		environment = os.Args[1]
	}
	
	// Set environment variable
	os.Setenv("PHONIC_ENV", environment)
	
	fmt.Printf("Testing health checks for environment: %s\n\n", environment)
	
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Initialize logger
	if err := logger.InitGlobal(cfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger := logger.GetGlobal()
	
	// Create health manager
	healthManager := health.NewManager(cfg.App.Name, cfg.App.Version, appLogger)
	
	fmt.Println("🔍 Setting up health checkers...")
	
	// Add database health checker (if available)
	if cfg.Database.Host != "" {
		db, err := setupDatabase(cfg)
		if err != nil {
			fmt.Printf("⚠️  Database setup failed: %v\n", err)
		} else {
			healthManager.AddChecker("database", health.NewDatabaseChecker(db, appLogger))
			fmt.Println("✅ Database health checker added")
		}
	}
	
	// Add Redis health checker (if available)
	if cfg.Redis.Host != "" {
		redisClient := setupRedis(cfg)
		healthManager.AddChecker("redis", health.NewRedisChecker(redisClient, appLogger))
		fmt.Println("✅ Redis health checker added")
	}
	
	// Add Moshi STT health checker
	if cfg.Moshi.STT.Host != "" {
		sttURL := fmt.Sprintf("%s:%d", cfg.Moshi.STT.Host, cfg.Moshi.STT.Port)
		healthManager.AddChecker("moshi_stt", health.NewMoshiChecker(sttURL, 5*time.Second, appLogger))
		fmt.Println("✅ Moshi STT health checker added")
	}
	
	// Add Moshi TTS health checker
	if cfg.Moshi.TTS.Host != "" {
		ttsURL := fmt.Sprintf("%s:%d", cfg.Moshi.TTS.Host, cfg.Moshi.TTS.Port)
		healthManager.AddChecker("moshi_tts", health.NewMoshiChecker(ttsURL, 5*time.Second, appLogger))
		fmt.Println("✅ Moshi TTS health checker added")
	}
	
	// Add custom health checker example
	healthManager.AddChecker("custom_check", health.NewCustomChecker(
		"system_resources",
		func(ctx context.Context) (bool, string, map[string]string) {
			// Simple example: check if we have enough disk space
			metadata := map[string]string{
				"check_type": "system_resources",
				"timestamp": time.Now().Format(time.RFC3339),
			}
			return true, "System resources OK", metadata
		},
		appLogger,
	))
	fmt.Println("✅ Custom health checker added")
	
	fmt.Println("\n🏥 Running health checks...")
	
	// Perform health check
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	healthResponse := healthManager.CheckHealth(ctx)
	
	// Display results
	fmt.Printf("\n📊 Health Check Results:\n")
	fmt.Printf("Service: %s\n", healthResponse.Service)
	fmt.Printf("Version: %s\n", healthResponse.Version)
	fmt.Printf("Overall Status: %s\n", healthResponse.Status)
	fmt.Printf("Uptime: %v\n", healthResponse.Uptime)
	fmt.Printf("Timestamp: %v\n", healthResponse.Timestamp.Format(time.RFC3339))
	fmt.Printf("Checks: %d\n\n", len(healthResponse.Checks))
	
	// Display individual check results
	for _, check := range healthResponse.Checks {
		status := "✅"
		if check.Status != health.StatusHealthy {
			status = "❌"
		}
		
		fmt.Printf("%s %s: %s (Duration: %v)\n", status, check.Name, check.Status, check.Duration)
		if check.Message != "" {
			fmt.Printf("   Message: %s\n", check.Message)
		}
		if len(check.Metadata) > 0 {
			fmt.Printf("   Metadata:\n")
			for key, value := range check.Metadata {
				fmt.Printf("     %s: %s\n", key, value)
			}
		}
		fmt.Println()
	}
	
	// Test HTTP endpoints
	fmt.Println("🌐 Testing HTTP health endpoints...")
	
	// Start a temporary HTTP server to test endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthManager.HTTPHandler())
	mux.HandleFunc("/health/ready", healthManager.ReadinessHandler())
	mux.HandleFunc("/health/live", healthManager.LivenessHandler())
	
	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}
	
	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Test endpoints
	testEndpoints := []string{
		"http://localhost:8888/health",
		"http://localhost:8888/health/ready",
		"http://localhost:8888/health/live",
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	
	for _, endpoint := range testEndpoints {
		resp, err := client.Get(endpoint)
		if err != nil {
			fmt.Printf("❌ %s: Failed - %v\n", endpoint, err)
			continue
		}
		
		status := "✅"
		if resp.StatusCode >= 400 {
			status = "❌"
		}
		
		fmt.Printf("%s %s: %d %s\n", status, endpoint, resp.StatusCode, resp.Status)
		resp.Body.Close()
	}
	
	// Shutdown server
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	
	fmt.Println("\n✅ Health check test completed!")
}

func setupDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

func setupRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
		PoolSize: cfg.Redis.PoolSize,
	})
	
	return client
}
