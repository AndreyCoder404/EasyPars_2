package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration structure
// This struct maps to the config.yaml file structure
type Config struct {
	// Server configuration section
	Server ServerConfig `mapstructure:"server" yaml:"server"`

	// Future configuration sections to be added:
	// Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	// Parser   ParserConfig   `mapstructure:"parser" yaml:"parser"`
	// JWT      JWTConfig      `mapstructure:"jwt" yaml:"jwt"`
	// Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	// Logging  LoggingConfig  `mapstructure:"logging" yaml:"logging"`
}

// ServerConfig holds server-specific configuration
// Maps to the "server" section in config.yaml
type ServerConfig struct {
	Port string `mapstructure:"port" yaml:"port"`

	// Future server configuration fields:
	// Host         string `mapstructure:"host" yaml:"host"`
	// ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	// WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	// TLSEnabled   bool   `mapstructure:"tls_enabled" yaml:"tls_enabled"`
	// CertFile     string `mapstructure:"cert_file" yaml:"cert_file"`
	// KeyFile      string `mapstructure:"key_file" yaml:"key_file"`
}

// DatabaseConfig will hold database configuration (future implementation)
// type DatabaseConfig struct {
// 	Host     string `mapstructure:"host" yaml:"host"`
// 	Port     int    `mapstructure:"port" yaml:"port"`
// 	User     string `mapstructure:"user" yaml:"user"`
// 	Password string `mapstructure:"password" yaml:"password"`
// 	DBName   string `mapstructure:"dbname" yaml:"dbname"`
// 	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
// }

// ParserConfig will hold parser configuration (future implementation)
// type ParserConfig struct {
// 	BaseURL           string `mapstructure:"base_url" yaml:"base_url"`
// 	RateLimit         int    `mapstructure:"rate_limit" yaml:"rate_limit"`
// 	Timeout           int    `mapstructure:"timeout" yaml:"timeout"`
// 	ConcurrentWorkers int    `mapstructure:"concurrent_workers" yaml:"concurrent_workers"`
// 	RetryAttempts     int    `mapstructure:"retry_attempts" yaml:"retry_attempts"`
// }

// JWTConfig will hold JWT configuration (future implementation)
// type JWTConfig struct {
// 	Secret      string `mapstructure:"secret" yaml:"secret"`
// 	ExpireHours int    `mapstructure:"expire_hours" yaml:"expire_hours"`
// 	Issuer      string `mapstructure:"issuer" yaml:"issuer"`
// }

// LoadConfig loads configuration from config.yaml using Viper
// This function initializes Viper, sets up configuration sources, and loads the config
func LoadConfig() (*Config, error) {
	// Initialize a new Viper instance
	// Viper is a configuration solution for Go applications
	v := viper.New()

	// Set the configuration file name (without extension)
	// Viper will look for config.yaml, config.yml, config.json, etc.
	v.SetConfigName("config")

	// Set the configuration file type explicitly
	// This ensures Viper knows we're working with YAML
	v.SetConfigType("yaml")

	// Add configuration search paths
	// Viper will search for the config file in these directories
	v.AddConfigPath(".")               // Current directory
	v.AddConfigPath("./config")        // Config subdirectory
	v.AddConfigPath("$HOME/.easypars") // User home directory (future)

	// Enable environment variable support
	// This allows overriding config values with environment variables
	// Format: EASYPARS_SERVER_PORT will override server.port
	v.SetEnvPrefix("EASYPARS")
	v.AutomaticEnv()

	// Set default values for critical configurations
	// These defaults ensure the application can start even without a config file
	setDefaultValues(v)

	// Read the configuration file
	// This step loads the config.yaml file into Viper
	if err := v.ReadInConfig(); err != nil {
		// Handle different types of configuration errors
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found - use defaults and log a warning
			log.Printf("Warning: Config file not found, using default values")
		} else {
			// Config file found but there's an error reading it
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		// Successfully loaded config file
		log.Printf("Config file loaded: %s", v.ConfigFileUsed())
	}

	// Create a new Config instance to hold the loaded configuration
	var config Config

	// Unmarshal the configuration into our Config struct
	// This maps the YAML structure to our Go struct using mapstructure tags
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate the loaded configuration
	// This ensures all required fields are present and valid
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Log successful configuration loading
	log.Printf("Configuration loaded successfully - Server will start on port: %s", config.Server.Port)

	return &config, nil
}

// setDefaultValues sets default configuration values
// This ensures the application has sensible defaults even without a config file
func setDefaultValues(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")

	// Future default values to be added:
	// v.SetDefault("server.host", "localhost")
	// v.SetDefault("server.read_timeout", 30)
	// v.SetDefault("server.write_timeout", 30)
	// v.SetDefault("database.host", "localhost")
	// v.SetDefault("database.port", 5432)
	// v.SetDefault("parser.base_url", "https://vringe.com/results/")
	// v.SetDefault("parser.rate_limit", 5)
	// v.SetDefault("parser.timeout", 30)
	// v.SetDefault("parser.concurrent_workers", 3)
	// v.SetDefault("jwt.expire_hours", 24)
}

// validateConfig validates the loaded configuration
// This function checks that all required fields are present and valid
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate port format (should be numeric or :port format)
	if !isValidPort(config.Server.Port) {
		return fmt.Errorf("invalid server port format: %s", config.Server.Port)
	}

	// Future validation to be added:
	// - Database connection parameters
	// - Parser URL format validation
	// - JWT secret strength validation
	// - File path existence checks

	return nil
}

// isValidPort checks if the port string is in a valid format
// Accepts formats like "8080", ":8080", or "localhost:8080"
func isValidPort(port string) bool {
	// Basic validation - ensure port is not empty
	if port == "" {
		return false
	}

	// Future validation to be added:
	// - Numeric port range validation (1-65535)
	// - Host:port format validation
	// - Reserved port checks

	return true
}

// GetConfigPath returns the path to the configuration file
// This is useful for logging and debugging configuration issues
func GetConfigPath() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}

	// Return the expected config file path
	return filepath.Join(cwd, "config.yaml")
}

// ReloadConfig reloads the configuration from file
// This function will be useful for hot-reloading configuration changes
// Future implementation for production use
func ReloadConfig() (*Config, error) {
	// Future implementation:
	// - Watch for file changes
	// - Reload configuration without restarting the application
	// - Validate new configuration before applying
	// - Notify components of configuration changes

	return LoadConfig()
}

// Future functions to be implemented:
// - WatchConfigChanges() - for hot reloading
// - MergeConfigs(*Config, *Config) - for merging multiple config sources
// - ExportConfig(*Config) - for exporting current config to file
// - EncryptSensitiveFields(*Config) - for encrypting passwords/secrets
// - ValidateEnvironment() - for environment-specific validations
