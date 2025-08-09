// Config test utility for Phonic AI Calling Agent
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ArbajAnsari19/phonic/pkg/config"
)

func main() {
	fmt.Println("🎵 Phonic Configuration Test")
	fmt.Println("============================")
	
	// Get environment from command line or default to dev
	environment := "dev"
	if len(os.Args) > 1 {
		environment = os.Args[1]
	}
	
	// Set environment variable
	os.Setenv("PHONIC_ENV", environment)
	
	fmt.Printf("Testing configuration for environment: %s\n\n", environment)
	
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Display configuration summary
	fmt.Printf("App Information:\n")
	fmt.Printf("  Name: %s\n", cfg.App.Name)
	fmt.Printf("  Version: %s\n", cfg.App.Version)
	fmt.Printf("  Environment: %s\n", cfg.App.Environment)
	fmt.Printf("  Debug: %v\n", cfg.App.Debug)
	fmt.Printf("  Listen: %s:%d\n", cfg.App.Host, cfg.App.Port)
	fmt.Println()
	
	fmt.Printf("Database Configuration:\n")
	fmt.Printf("  Host: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("  Database: %s\n", cfg.Database.Database)
	fmt.Printf("  Username: %s\n", cfg.Database.Username)
	fmt.Printf("  SSL Mode: %s\n", cfg.Database.SSLMode)
	fmt.Printf("  Max Connections: %d/%d\n", cfg.Database.MaxIdleConns, cfg.Database.MaxOpenConns)
	fmt.Printf("  Connection URL: %s\n", cfg.GetDatabaseURL())
	fmt.Println()
	
	fmt.Printf("Redis Configuration:\n")
	fmt.Printf("  Address: %s\n", cfg.GetRedisAddr())
	fmt.Printf("  Database: %d\n", cfg.Redis.Database)
	fmt.Printf("  Pool Size: %d\n", cfg.Redis.PoolSize)
	fmt.Println()
	
	fmt.Printf("Moshi Configuration:\n")
	fmt.Printf("  STT URL: %s\n", cfg.GetMoshiSTTURL())
	fmt.Printf("  TTS URL: %s\n", cfg.GetMoshiTTSURL())
	fmt.Printf("  STT Sample Rate: %d Hz\n", cfg.Moshi.STT.SampleRate)
	fmt.Printf("  TTS Voice: %s (speed: %.1f)\n", cfg.Moshi.TTS.VoiceID, cfg.Moshi.TTS.Speed)
	fmt.Println()
	
	fmt.Printf("Storage Configuration:\n")
	fmt.Printf("  Endpoint: %s\n", cfg.Storage.Endpoint)
	fmt.Printf("  Bucket: %s\n", cfg.Storage.Bucket)
	fmt.Printf("  SSL: %v\n", cfg.Storage.UseSSL)
	fmt.Printf("  Retention: %d days\n", cfg.Storage.AudioRetentionDays)
	fmt.Println()
	
	fmt.Printf("Logging Configuration:\n")
	fmt.Printf("  Level: %s\n", cfg.Logging.Level)
	fmt.Printf("  Format: %s\n", cfg.Logging.Format)
	fmt.Printf("  Output: %s\n", cfg.Logging.Output)
	fmt.Println()
	
	fmt.Printf("Environment Flags:\n")
	fmt.Printf("  Is Development: %v\n", cfg.IsDevelopment())
	fmt.Printf("  Is Production: %v\n", cfg.IsProduction())
	fmt.Println()
	
	fmt.Println("✅ Configuration loaded and validated successfully!")
}
