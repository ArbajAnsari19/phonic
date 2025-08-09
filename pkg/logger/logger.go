// Package logger provides structured logging for Phonic AI Calling Agent
package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ArbajAnsari19/phonic/pkg/config"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
	config *config.Config
}

// Fields represents structured log fields
type Fields map[string]interface{}

// ContextKey is used for context values
type ContextKey string

const (
	// TraceIDKey is the context key for trace IDs
	TraceIDKey ContextKey = "traceID"
	// RequestIDKey is the context key for request IDs
	RequestIDKey ContextKey = "requestID"
	// ServiceKey is the context key for service name
	ServiceKey ContextKey = "service"
)

var (
	// Global logger instance
	globalLogger *Logger
)

// New creates a new logger instance based on configuration
func New(cfg *config.Config) (*Logger, error) {
	var zapConfig zap.Config
	
	// Configure based on environment
	if cfg.IsDevelopment() {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.MessageKey = "message"
		zapConfig.EncoderConfig.LevelKey = "level"
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.CallerKey = "caller"
	}
	
	// Set log level based on configuration
	level, err := zapcore.ParseLevel(cfg.Logging.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	
	// Configure output format
	switch cfg.Logging.Format {
	case "json":
		zapConfig.Encoding = "json"
	case "console":
		zapConfig.Encoding = "console"
	default:
		if cfg.IsDevelopment() {
			zapConfig.Encoding = "console"
		} else {
			zapConfig.Encoding = "json"
		}
	}
	
	// Configure output destination
	switch cfg.Logging.Output {
	case "stdout":
		zapConfig.OutputPaths = []string{"stdout"}
	case "stderr":
		zapConfig.OutputPaths = []string{"stderr"}
	default:
		zapConfig.OutputPaths = []string{"stdout"}
	}
	
	// Build the logger
	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}
	
	// Add caller information for debugging
	zapLogger = zapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	
	// Add service information
	zapLogger = zapLogger.With(
		zap.String("service", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Environment),
	)
	
	logger := &Logger{
		Logger: zapLogger,
		config: cfg,
	}
	
	return logger, nil
}

// InitGlobal initializes the global logger instance
func InitGlobal(cfg *config.Config) error {
	logger, err := New(cfg)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetGlobal returns the global logger instance
func GetGlobal() *Logger {
	if globalLogger == nil {
		// Fallback logger if global not initialized
		zapLogger, _ := zap.NewDevelopment()
		return &Logger{Logger: zapLogger}
	}
	return globalLogger
}

// WithContext creates a logger with context information
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := []zap.Field{}
	
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		fields = append(fields, zap.String("trace_id", traceID.(string)))
	}
	
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields = append(fields, zap.String("request_id", requestID.(string)))
	}
	
	if service := ctx.Value(ServiceKey); service != nil {
		fields = append(fields, zap.String("service_name", service.(string)))
	}
	
	return &Logger{
		Logger: l.Logger.With(fields...),
		config: l.config,
	}
}

// WithFields creates a logger with additional fields
func (l *Logger) WithFields(fields Fields) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	
	return &Logger{
		Logger: l.Logger.With(zapFields...),
		config: l.config,
	}
}

// WithService creates a logger with service name
func (l *Logger) WithService(serviceName string) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.String("service_name", serviceName)),
		config: l.config,
	}
}

// InfoWithDuration logs info with duration measurement
func (l *Logger) InfoWithDuration(msg string, duration time.Duration, fields ...zap.Field) {
	allFields := append(fields, zap.Duration("duration", duration))
	l.Info(msg, allFields...)
}

// ErrorWithDuration logs error with duration measurement
func (l *Logger) ErrorWithDuration(msg string, duration time.Duration, err error, fields ...zap.Field) {
	allFields := append(fields, 
		zap.Duration("duration", duration),
		zap.Error(err),
	)
	l.Error(msg, allFields...)
}

// LogHTTPRequest logs HTTP request information
func (l *Logger) LogHTTPRequest(method, path, userAgent, remoteAddr string, statusCode int, duration time.Duration) {
	l.Info("HTTP request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("user_agent", userAgent),
		zap.String("remote_addr", remoteAddr),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	)
}

// LogGRPCRequest logs gRPC request information
func (l *Logger) LogGRPCRequest(method string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("grpc_method", method),
		zap.Duration("duration", duration),
	}
	
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("gRPC request failed", fields...)
	} else {
		l.Info("gRPC request", fields...)
	}
}

// LogDatabaseOperation logs database operation
func (l *Logger) LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("db_operation", operation),
		zap.String("db_table", table),
		zap.Duration("duration", duration),
	}
	
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Database operation failed", fields...)
	} else {
		l.Info("Database operation", fields...)
	}
}

// LogWebSocketEvent logs WebSocket events
func (l *Logger) LogWebSocketEvent(event, connectionID string, fields Fields) {
	zapFields := []zap.Field{
		zap.String("ws_event", event),
		zap.String("connection_id", connectionID),
	}
	
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	
	l.Info("WebSocket event", zapFields...)
}

// LogAudioProcessing logs audio processing events
func (l *Logger) LogAudioProcessing(operation string, duration time.Duration, audioLength time.Duration, sampleRate int, channels int) {
	l.Info("Audio processing",
		zap.String("audio_operation", operation),
		zap.Duration("processing_duration", duration),
		zap.Duration("audio_length", audioLength),
		zap.Int("sample_rate", sampleRate),
		zap.Int("channels", channels),
	)
}

// LogMoshiInteraction logs interactions with Moshi STT/TTS services
func (l *Logger) LogMoshiInteraction(service, operation string, duration time.Duration, success bool, err error) {
	fields := []zap.Field{
		zap.String("moshi_service", service),
		zap.String("moshi_operation", operation),
		zap.Duration("duration", duration),
		zap.Bool("success", success),
	}
	
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Moshi interaction failed", fields...)
	} else {
		l.Info("Moshi interaction", fields...)
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Cleanup properly closes the logger
func (l *Logger) Cleanup() error {
	return l.Logger.Sync()
}

// Global convenience functions

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...zap.Field) {
	GetGlobal().Debug(msg, fields...)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...zap.Field) {
	GetGlobal().Info(msg, fields...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...zap.Field) {
	GetGlobal().Warn(msg, fields...)
}

// Error logs an error message using the global logger
func Error(msg string, fields ...zap.Field) {
	GetGlobal().Error(msg, fields...)
}

// Fatal logs a fatal message and exits using the global logger
func Fatal(msg string, fields ...zap.Field) {
	GetGlobal().Fatal(msg, fields...)
}

// WithContext creates a logger with context using the global logger
func WithContext(ctx context.Context) *Logger {
	return GetGlobal().WithContext(ctx)
}

// WithFields creates a logger with fields using the global logger
func WithFields(fields Fields) *Logger {
	return GetGlobal().WithFields(fields)
}

// WithService creates a logger with service name using the global logger
func WithService(serviceName string) *Logger {
	return GetGlobal().WithService(serviceName)
}

// CreateStartupLogger creates a basic logger for application startup
func CreateStartupLogger() *Logger {
	// Create a simple logger for startup before config is loaded
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLogger, err := zapConfig.Build()
	if err != nil {
		// Fallback to a no-op logger
		zapLogger = zap.NewNop()
	}
	
	return &Logger{Logger: zapLogger}
}

// SetGlobalLevel dynamically sets the global log level
func SetGlobalLevel(level string) error {
	if globalLogger == nil {
		return fmt.Errorf("global logger not initialized")
	}
	
	_, err := zapcore.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	
	// Note: This would require a reconfigurable logger setup
	// For now, we'll just log the change
	globalLogger.Info("Log level change requested", 
		zap.String("new_level", level),
		zap.String("note", "restart required to apply"),
	)
	
	return nil
}
