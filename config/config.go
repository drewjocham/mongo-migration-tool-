// Package config provides configuration management for mongo-essential.
//
// This package handles loading configuration from environment variables and .env files,
// supporting MongoDB connections, SSL/TLS settings, AI provider configuration, and
// Google Docs integration.
//
// Basic Usage:
//
//	// Load configuration from .env files
//	cfg, err := config.Load(".env", ".env.local")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Use configuration to connect to MongoDB
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetConnectionString()))
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Environment Variables:
//
// MongoDB Configuration:
//   - MONGO_URL: MongoDB connection URL
//   - MONGO_DATABASE: Database name (required)
//   - MONGO_USERNAME: Authentication username
//   - MONGO_PASSWORD: Authentication password
//
// SSL/TLS Settings:
//   - MONGO_SSL_ENABLED: Enable SSL/TLS connections
//   - MONGO_SSL_INSECURE: Skip certificate verification (useful for development)
//
// AI Analysis:
//   - AI_ENABLED: Enable AI-powered database analysis
//   - AI_PROVIDER: AI provider (openai, gemini, claude)
//   - OPENAI_API_KEY: OpenAI API key
//   - GEMINI_API_KEY: Google Gemini API key
//   - CLAUDE_API_KEY: Anthropic Claude API key
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds all configuration options for mongo-essential.
//
// Configuration is loaded from environment variables with support for .env files.
// All MongoDB connection settings, AI provider configurations, and Google Docs
// integration settings are contained within this struct.
type Config struct {
	// MongoDB connection settings
	MongoURL string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
	Database string `env:"MONGO_DATABASE,required"`

	// Migration settings
	MigrationsPath       string `env:"MIGRATIONS_PATH" envDefault:"./migrations"`
	MigrationsCollection string `env:"MIGRATIONS_COLLECTION" envDefault:"schema_migrations"`

	// Authentication
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`

	// SSL/TLS settings - important for cloud providers like STACKIT
	SSLEnabled           bool   `env:"MONGO_SSL_ENABLED" envDefault:"false"`
	SSLInsecure          bool   `env:"MONGO_SSL_INSECURE" envDefault:"false"`
	SSLCertificatePath   string `env:"MONGO_SSL_CERT_PATH"`
	SSLPrivateKeyPath    string `env:"MONGO_SSL_KEY_PATH"`
	SSLCACertificatePath string `env:"MONGO_SSL_CA_CERT_PATH"`

	// Connection pool settings
	MaxPoolSize int `env:"MONGO_MAX_POOL_SIZE" envDefault:"10"`
	MinPoolSize int `env:"MONGO_MIN_POOL_SIZE" envDefault:"1"`
	MaxIdleTime int `env:"MONGO_MAX_IDLE_TIME" envDefault:"300"` // in seconds

	// Connection timeout
	Timeout int `env:"MONGO_TIMEOUT" envDefault:"60"`

	// AI Analysis settings
	AIProvider string `env:"AI_PROVIDER" envDefault:"openai"`
	AIEnabled  bool   `env:"AI_ENABLED" envDefault:"false"`

	// OpenAI settings
	OpenAIAPIKey string `env:"OPENAI_API_KEY"`
	OpenAIModel  string `env:"OPENAI_MODEL" envDefault:"gpt-4o-mini"`

	// Google Gemini settings
	GeminiAPIKey string `env:"GEMINI_API_KEY"`
	GeminiModel  string `env:"GEMINI_MODEL" envDefault:"gemini-1.5-flash"`

	// Claude settings
	ClaudeAPIKey string `env:"CLAUDE_API_KEY"`
	ClaudeModel  string `env:"CLAUDE_MODEL" envDefault:"claude-3-5-sonnet-20241022"`

	// Google Docs Integration
	GoogleDocsEnabled        bool   `env:"GOOGLE_DOCS_ENABLED" envDefault:"false"`
	GoogleCredentialsPath    string `env:"GOOGLE_CREDENTIALS_PATH"`
	GoogleCredentialsJSON    string `env:"GOOGLE_CREDENTIALS_JSON"`
	GoogleDriveFolderID      string `env:"GOOGLE_DRIVE_FOLDER_ID"`
	GoogleDocsTemplate       string `env:"GOOGLE_DOCS_TEMPLATE" envDefault:"analysis"`
	GoogleDocsShareWithEmail string `env:"GOOGLE_DOCS_SHARE_WITH_EMAIL"`
}

// Load loads configuration from the specified .env files and environment variables.
//
// Files are loaded in order, with later files overriding earlier ones.
// After loading files, environment variables take precedence over file values.
//
// Example:
//
//	// Load from .env and .env.local
//	cfg, err := config.Load(".env", ".env.local")
//
//	// Load from custom file
//	cfg, err := config.Load("config/production.env")
//
// If a file doesn't exist, it will be silently skipped.
func Load(envFiles ...string) (*Config, error) {
	for _, file := range envFiles {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err != nil {
				return nil, fmt.Errorf("failed to load env file %s: %w", file, err)
			}
		}
	}

	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

// LoadFromEnv loads configuration only from environment variables, ignoring .env files.
//
// This is useful in containerized environments where configuration is provided
// through environment variables rather than files.
//
// Example:
//
//	cfg, err := config.LoadFromEnv()
func LoadFromEnv() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

// GetConnectionString builds the MongoDB connection string with authentication.
//
// This method handles various MongoDB connection scenarios:
// - Local development (mongodb://localhost:27017)
// - Cloud providers with authentication (mongodb+srv://user:pass@cluster.provider.com)
// - Custom connection strings with credential injection
//
// The method automatically injects username/password credentials if they are
// provided and not already present in the connection URL.
//
// Example:
//
//	cfg := &Config{
//		MongoURL: "mongodb+srv://cluster.example.com/database",
//		Username: "user",
//		Password: "pass",
//	}
//	// Returns: "mongodb+srv://user:pass@cluster.example.com/database"
//	connectionString := cfg.GetConnectionString()
func (c *Config) GetConnectionString() string {
	connStr := c.MongoURL

	if c.Username != "" && c.Password != "" {
		if c.MongoURL == "mongodb://localhost:27017" ||
			(c.MongoURL != "" && !strings.Contains(c.MongoURL, "@")) {
			if strings.HasPrefix(c.MongoURL, "mongodb://") {
				// Replace mongodb:// with mongodb://user:pass@
				connStr = strings.Replace(c.MongoURL, "mongodb://",
					fmt.Sprintf("mongodb://%s:%s@", c.Username, c.Password), 1)
			} else if strings.HasPrefix(c.MongoURL, "mongodb+srv://") {
				// Replace mongodb+srv:// with mongodb+srv://user:pass@
				connStr = strings.Replace(c.MongoURL, "mongodb+srv://",
					fmt.Sprintf("mongodb+srv://%s:%s@", c.Username, c.Password), 1)
			}
		}
	}

	return connStr
}

// Validate performs basic validation on the configuration.
//
// This checks for required fields and validates that AI provider settings
// are consistent (e.g., if AI is enabled, appropriate API keys are provided).
//
// Returns an error if validation fails.
func (c *Config) Validate() error {
	if c.Database == "" {
		return fmt.Errorf("MONGO_DATABASE is required")
	}

	if c.AIEnabled {
		switch c.AIProvider {
		case "openai":
			if c.OpenAIAPIKey == "" {
				return fmt.Errorf("OPENAI_API_KEY is required when AI_PROVIDER is openai")
			}
		case "gemini":
			if c.GeminiAPIKey == "" {
				return fmt.Errorf("GEMINI_API_KEY is required when AI_PROVIDER is gemini")
			}
		case "claude":
			if c.ClaudeAPIKey == "" {
				return fmt.Errorf("CLAUDE_API_KEY is required when AI_PROVIDER is claude")
			}
		default:
			return fmt.Errorf("invalid AI_PROVIDER: %s (must be openai, gemini, or claude)", c.AIProvider)
		}
	}

	if c.GoogleDocsEnabled {
		if c.GoogleCredentialsPath == "" && c.GoogleCredentialsJSON == "" {
			return fmt.Errorf("either GOOGLE_CREDENTIALS_PATH or GOOGLE_CREDENTIALS_JSON is required when GOOGLE_DOCS_ENABLED is true")
		}
	}

	return nil
}
