package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	MongoURL string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
	Database string `env:"MONGO_DATABASE,required"`

	MigrationsPath       string `env:"MIGRATIONS_PATH" envDefault:"./migrations"`
	MigrationsCollection string `env:"MIGRATIONS_COLLECTION" envDefault:"schema_migrations"`

	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`

	// SSL/TLS settings
	SSLEnabled           bool   `env:"MONGO_SSL_ENABLED" envDefault:"false"`
	SSLInsecure          bool   `env:"MONGO_SSL_INSECURE" envDefault:"false"`
	SSLCertificatePath   string `env:"MONGO_SSL_CERT_PATH"`
	SSLPrivateKeyPath    string `env:"MONGO_SSL_KEY_PATH"`
	SSLCACertificatePath string `env:"MONGO_SSL_CA_CERT_PATH"`

	MaxPoolSize int `env:"MONGO_MAX_POOL_SIZE" envDefault:"10"`
	MinPoolSize int `env:"MONGO_MIN_POOL_SIZE" envDefault:"1"`
	MaxIdleTime int `env:"MONGO_MAX_IDLE_TIME" envDefault:"300"` // in seconds

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

// GetConnectionString builds the MongoDB connection string with cloud provider support.
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
