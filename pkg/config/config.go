// Package config provides centralized configuration management for Phonic AI Calling Agent
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	Moshi    MoshiConfig    `mapstructure:"moshi" yaml:"moshi"`
	Services ServicesConfig `mapstructure:"services" yaml:"services"`
	Logging  LoggingConfig  `mapstructure:"logging" yaml:"logging"`
	Security SecurityConfig `mapstructure:"security" yaml:"security"`
	Storage  StorageConfig  `mapstructure:"storage" yaml:"storage"`
}

// AppConfig contains general application settings
type AppConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Version     string `mapstructure:"version" yaml:"version"`
	Environment string `mapstructure:"environment" yaml:"environment"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
	Port        int    `mapstructure:"port" yaml:"port"`
	Host        string `mapstructure:"host" yaml:"host"`
}

// DatabaseConfig contains PostgreSQL connection settings
type DatabaseConfig struct {
	Host            string        `mapstructure:"host" yaml:"host"`
	Port            int           `mapstructure:"port" yaml:"port"`
	Username        string        `mapstructure:"username" yaml:"username"`
	Password        string        `mapstructure:"password" yaml:"password"`
	Database        string        `mapstructure:"database" yaml:"database"`
	SSLMode         string        `mapstructure:"ssl_mode" yaml:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	Database int    `mapstructure:"database" yaml:"database"`
	PoolSize int    `mapstructure:"pool_size" yaml:"pool_size"`
}

// MoshiConfig contains Kyutai Moshi server settings
type MoshiConfig struct {
	STT MoshiSTTConfig `mapstructure:"stt" yaml:"stt"`
	TTS MoshiTTSConfig `mapstructure:"tts" yaml:"tts"`
}

// MoshiSTTConfig contains STT server settings
type MoshiSTTConfig struct {
	Host           string        `mapstructure:"host" yaml:"host"`
	Port           int           `mapstructure:"port" yaml:"port"`
	WebSocketPath  string        `mapstructure:"websocket_path" yaml:"websocket_path"`
	Timeout        time.Duration `mapstructure:"timeout" yaml:"timeout"`
	RetryAttempts  int           `mapstructure:"retry_attempts" yaml:"retry_attempts"`
	SampleRate     int           `mapstructure:"sample_rate" yaml:"sample_rate"`
	Channels       int           `mapstructure:"channels" yaml:"channels"`
	ChunkSize      int           `mapstructure:"chunk_size" yaml:"chunk_size"`
}

// MoshiTTSConfig contains TTS server settings
type MoshiTTSConfig struct {
	Host          string        `mapstructure:"host" yaml:"host"`
	Port          int           `mapstructure:"port" yaml:"port"`
	WebSocketPath string        `mapstructure:"websocket_path" yaml:"websocket_path"`
	Timeout       time.Duration `mapstructure:"timeout" yaml:"timeout"`
	RetryAttempts int           `mapstructure:"retry_attempts" yaml:"retry_attempts"`
	VoiceID       string        `mapstructure:"voice_id" yaml:"voice_id"`
	Speed         float64       `mapstructure:"speed" yaml:"speed"`
	Quality       string        `mapstructure:"quality" yaml:"quality"`
}

// ServicesConfig contains settings for other microservices
type ServicesConfig struct {
	Gateway      ServiceEndpoint `mapstructure:"gateway" yaml:"gateway"`
	Session      ServiceEndpoint `mapstructure:"session" yaml:"session"`
	Orchestrator ServiceEndpoint `mapstructure:"orchestrator" yaml:"orchestrator"`
	STTClient    ServiceEndpoint `mapstructure:"stt_client" yaml:"stt_client"`
	TTSClient    ServiceEndpoint `mapstructure:"tts_client" yaml:"tts_client"`
}

// ServiceEndpoint represents a microservice endpoint
type ServiceEndpoint struct {
	Host    string        `mapstructure:"host" yaml:"host"`
	Port    int           `mapstructure:"port" yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
	Output string `mapstructure:"output" yaml:"output"`
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	JWTSecret    string        `mapstructure:"jwt_secret" yaml:"jwt_secret"`
	JWTExpiryHours int         `mapstructure:"jwt_expiry_hours" yaml:"jwt_expiry_hours"`
	RateLimit    RateLimitConfig `mapstructure:"rate_limit" yaml:"rate_limit"`
	CORS         CORSConfig    `mapstructure:"cors" yaml:"cors"`
}

// RateLimitConfig contains rate limiting settings
type RateLimitConfig struct {
	RequestsPerMinute int           `mapstructure:"requests_per_minute" yaml:"requests_per_minute"`
	BurstSize         int           `mapstructure:"burst_size" yaml:"burst_size"`
	WindowSize        time.Duration `mapstructure:"window_size" yaml:"window_size"`
}

// CORSConfig contains CORS settings
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers" yaml:"allowed_headers"`
}

// StorageConfig contains MinIO/S3 settings
type StorageConfig struct {
	Endpoint        string `mapstructure:"endpoint" yaml:"endpoint"`
	AccessKey       string `mapstructure:"access_key" yaml:"access_key"`
	SecretKey       string `mapstructure:"secret_key" yaml:"secret_key"`
	Bucket          string `mapstructure:"bucket" yaml:"bucket"`
	Region          string `mapstructure:"region" yaml:"region"`
	UseSSL          bool   `mapstructure:"use_ssl" yaml:"use_ssl"`
	AudioRetentionDays int `mapstructure:"audio_retention_days" yaml:"audio_retention_days"`
}

// Load loads configuration from files and environment variables
func Load(configPath string) (*Config, error) {
	// Set default configuration file name and paths
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	
	// Add configuration paths
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	
	// Default config paths based on environment
	environment := getEnvironment()
	viper.AddConfigPath(fmt.Sprintf("./configs/%s", environment))
	viper.AddConfigPath("./configs/dev")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	
	// Environment variable configuration
	viper.SetEnvPrefix("PHONIC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	
	// Set defaults
	setDefaults()
	
	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; using defaults and env vars
			fmt.Printf("Warning: Config file not found, using defaults and environment variables\n")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &config, nil
}

// getEnvironment returns the current environment
func getEnvironment() string {
	env := os.Getenv("PHONIC_ENV")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "dev"
	}
	return env
}

// setDefaults sets default configuration values
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "Phonic AI Calling Agent")
	viper.SetDefault("app.environment", "dev")
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.host", "0.0.0.0")
	
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "phonic")
	viper.SetDefault("database.password", "phonic_dev_password")
	viper.SetDefault("database.database", "phonic")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", "1h")
	
	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.pool_size", 10)
	
	// Moshi defaults
	viper.SetDefault("moshi.stt.host", "localhost")
	viper.SetDefault("moshi.stt.port", 8001)
	viper.SetDefault("moshi.stt.websocket_path", "/transcribe")
	viper.SetDefault("moshi.stt.timeout", "30s")
	viper.SetDefault("moshi.stt.retry_attempts", 3)
	viper.SetDefault("moshi.stt.sample_rate", 16000)
	viper.SetDefault("moshi.stt.channels", 1)
	viper.SetDefault("moshi.stt.chunk_size", 1600)
	
	viper.SetDefault("moshi.tts.host", "localhost")
	viper.SetDefault("moshi.tts.port", 8002)
	viper.SetDefault("moshi.tts.websocket_path", "/synthesize")
	viper.SetDefault("moshi.tts.timeout", "30s")
	viper.SetDefault("moshi.tts.retry_attempts", 3)
	viper.SetDefault("moshi.tts.voice_id", "default")
	viper.SetDefault("moshi.tts.speed", 1.0)
	viper.SetDefault("moshi.tts.quality", "high")
	
	// Service defaults
	viper.SetDefault("services.gateway.host", "localhost")
	viper.SetDefault("services.gateway.port", 8080)
	viper.SetDefault("services.gateway.timeout", "30s")
	
	viper.SetDefault("services.session.host", "localhost")
	viper.SetDefault("services.session.port", 8083)
	viper.SetDefault("services.session.timeout", "30s")
	
	viper.SetDefault("services.orchestrator.host", "localhost")
	viper.SetDefault("services.orchestrator.port", 8084)
	viper.SetDefault("services.orchestrator.timeout", "30s")
	
	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	
	// Security defaults
	viper.SetDefault("security.jwt_expiry_hours", 24)
	viper.SetDefault("security.rate_limit.requests_per_minute", 100)
	viper.SetDefault("security.rate_limit.burst_size", 50)
	viper.SetDefault("security.rate_limit.window_size", "1m")
	
	// Storage defaults
	viper.SetDefault("storage.endpoint", "localhost:9000")
	viper.SetDefault("storage.access_key", "phonic")
	viper.SetDefault("storage.secret_key", "phonic_dev_password")
	viper.SetDefault("storage.bucket", "phonic-audio")
	viper.SetDefault("storage.region", "us-east-1")
	viper.SetDefault("storage.use_ssl", false)
	viper.SetDefault("storage.audio_retention_days", 30)
}

// validateConfig validates the loaded configuration
func validateConfig(config *Config) error {
	// Validate required fields
	if config.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	
	if config.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	
	if config.Redis.Host == "" {
		return fmt.Errorf("redis.host is required")
	}
	
	if config.Moshi.STT.Host == "" {
		return fmt.Errorf("moshi.stt.host is required")
	}
	
	if config.Moshi.TTS.Host == "" {
		return fmt.Errorf("moshi.tts.host is required")
	}
	
	// Validate environment
	validEnvs := []string{"dev", "staging", "prod"}
	envValid := false
	for _, env := range validEnvs {
		if config.App.Environment == env {
			envValid = true
			break
		}
	}
	if !envValid {
		return fmt.Errorf("app.environment must be one of: %v", validEnvs)
	}
	
	// Validate logging level
	validLevels := []string{"debug", "info", "warn", "error"}
	levelValid := false
	for _, level := range validLevels {
		if config.Logging.Level == level {
			levelValid = true
			break
		}
	}
	if !levelValid {
		return fmt.Errorf("logging.level must be one of: %v", validLevels)
	}
	
	return nil
}

// GetDatabaseURL returns a formatted database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns Redis address in host:port format
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetMoshiSTTURL returns the Moshi STT WebSocket URL
func (c *Config) GetMoshiSTTURL() string {
	protocol := "ws"
	return fmt.Sprintf("%s://%s:%d%s", protocol, c.Moshi.STT.Host, c.Moshi.STT.Port, c.Moshi.STT.WebSocketPath)
}

// GetMoshiTTSURL returns the Moshi TTS WebSocket URL
func (c *Config) GetMoshiTTSURL() string {
	protocol := "ws"
	return fmt.Sprintf("%s://%s:%d%s", protocol, c.Moshi.TTS.Host, c.Moshi.TTS.Port, c.Moshi.TTS.WebSocketPath)
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "dev"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Environment == "prod"
}
