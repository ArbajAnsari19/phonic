// Logging test utility for Phonic AI Calling Agent
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/ArbajAnsari19/phonic/pkg/config"
	"github.com/ArbajAnsari19/phonic/pkg/logger"
)

func main() {
	fmt.Println("🎵 Phonic Logging Test")
	fmt.Println("======================")
	
	// Get environment from command line or default to dev
	environment := "dev"
	if len(os.Args) > 1 {
		environment = os.Args[1]
	}
	
	// Set environment variable
	os.Setenv("PHONIC_ENV", environment)
	
	fmt.Printf("Testing logging for environment: %s\n\n", environment)
	
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Initialize logger
	appLogger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer appLogger.Cleanup()
	
	// Initialize global logger
	if err := logger.InitGlobal(cfg); err != nil {
		log.Fatalf("Failed to initialize global logger: %v", err)
	}
	
	fmt.Printf("Logger initialized for %s environment\n", cfg.App.Environment)
	fmt.Printf("Log Level: %s, Format: %s\n\n", cfg.Logging.Level, cfg.Logging.Format)
	
	// Test basic logging levels
	fmt.Println("Testing basic log levels:")
	appLogger.Debug("This is a debug message", zap.String("test", "debug"))
	appLogger.Info("This is an info message", zap.String("test", "info"))
	appLogger.Warn("This is a warning message", zap.String("test", "warn"))
	appLogger.Error("This is an error message", zap.String("test", "error"))
	
	// Test structured logging
	fmt.Println("\nTesting structured logging:")
	appLogger.Info("User login attempt",
		zap.String("user_id", "user123"),
		zap.String("email", "test@example.com"),
		zap.String("ip_address", "192.168.1.1"),
		zap.Bool("success", true),
		zap.Duration("response_time", 150*time.Millisecond),
	)
	
	// Test context-based logging
	fmt.Println("\nTesting context-based logging:")
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.TraceIDKey, "trace-12345")
	ctx = context.WithValue(ctx, logger.RequestIDKey, "req-67890")
	ctx = context.WithValue(ctx, logger.ServiceKey, "gateway")
	
	contextLogger := appLogger.WithContext(ctx)
	contextLogger.Info("Processing user request",
		zap.String("action", "process_audio"),
		zap.Int("user_id", 12345),
	)
	
	// Test service-specific logging
	fmt.Println("\nTesting service-specific logging:")
	gatewayLogger := appLogger.WithService("gateway")
	gatewayLogger.Info("WebRTC connection established",
		zap.String("connection_id", "conn-12345"),
		zap.String("client_ip", "10.0.0.1"),
	)
	
	// Test fields-based logging
	fmt.Println("\nTesting fields-based logging:")
	fieldsLogger := appLogger.WithFields(logger.Fields{
		"module": "audio_processor",
		"version": "1.0.0",
		"worker_id": 5,
	})
	fieldsLogger.Info("Audio processing started",
		zap.String("audio_format", "wav"),
		zap.Int("sample_rate", 16000),
	)
	
	// Test specialized logging methods
	fmt.Println("\nTesting specialized logging methods:")
	
	// HTTP request logging
	appLogger.LogHTTPRequest("POST", "/api/calls", "Mozilla/5.0", "192.168.1.100", 200, 250*time.Millisecond)
	
	// gRPC request logging
	appLogger.LogGRPCRequest("phonic.Session/CreateSession", 100*time.Millisecond, nil)
	appLogger.LogGRPCRequest("phonic.Session/GetSession", 50*time.Millisecond, fmt.Errorf("session not found"))
	
	// Database operation logging
	appLogger.LogDatabaseOperation("INSERT", "call_sessions", 25*time.Millisecond, nil)
	appLogger.LogDatabaseOperation("SELECT", "users", 15*time.Millisecond, fmt.Errorf("connection timeout"))
	
	// WebSocket event logging
	appLogger.LogWebSocketEvent("connection_opened", "ws-conn-123", logger.Fields{
		"client_ip": "10.0.0.5",
		"user_id": 456,
	})
	
	// Audio processing logging
	appLogger.LogAudioProcessing("transcription", 500*time.Millisecond, 5*time.Second, 16000, 1)
	
	// Moshi interaction logging
	appLogger.LogMoshiInteraction("stt", "transcribe", 200*time.Millisecond, true, nil)
	appLogger.LogMoshiInteraction("tts", "synthesize", 800*time.Millisecond, false, fmt.Errorf("connection refused"))
	
	// Test global logger functions
	fmt.Println("\nTesting global logger functions:")
	logger.Info("Global logger info message", zap.String("component", "test"))
	logger.Error("Global logger error message", zap.String("component", "test"), zap.Error(fmt.Errorf("test error")))
	
	// Test duration logging
	fmt.Println("\nTesting duration logging:")
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // Simulate work
	duration := time.Since(start)
	
	appLogger.InfoWithDuration("Operation completed successfully", duration,
		zap.String("operation", "test_operation"),
		zap.Int("records_processed", 100),
	)
	
	appLogger.ErrorWithDuration("Operation failed", duration, fmt.Errorf("something went wrong"),
		zap.String("operation", "test_operation"),
		zap.Int("records_processed", 50),
	)
	
	fmt.Println("\n✅ Logging test completed successfully!")
	
	// Test log level change
	fmt.Println("\nTesting dynamic log level change:")
	if err := logger.SetGlobalLevel("error"); err != nil {
		appLogger.Error("Failed to change log level", zap.Error(err))
	}
	
	// These should not appear if level change worked (but it requires restart)
	logger.Info("This info message should be suppressed if level changed to error")
	logger.Error("This error message should still appear")
}
