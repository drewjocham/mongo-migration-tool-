package config

import (
	"os"
	"testing"
)

func TestConfig_GetConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "basic connection string without credentials",
			config: Config{
				MongoURL: "mongodb://localhost:27017",
			},
			expected: "mongodb://localhost:27017",
		},
		{
			name: "mongodb:// with credentials",
			config: Config{
				MongoURL: "mongodb://localhost:27017",
				Username: "user",
				Password: "pass",
			},
			expected: "mongodb://user:pass@localhost:27017",
		},
		{
			name: "mongodb+srv:// with credentials",
			config: Config{
				MongoURL: "mongodb+srv://cluster.example.com/database",
				Username: "user",
				Password: "pass",
			},
			expected: "mongodb+srv://user:pass@cluster.example.com/database",
		},
		{
			name: "connection string already contains credentials",
			config: Config{
				MongoURL: "mongodb://existinguser:existingpass@cluster.example.com/database",
				Username: "user",
				Password: "pass",
			},
			expected: "mongodb://existinguser:existingpass@cluster.example.com/database",
		},
		{
			name: "credentials without password",
			config: Config{
				MongoURL: "mongodb://localhost:27017",
				Username: "user",
			},
			expected: "mongodb://localhost:27017",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetConnectionString()
			if result != tt.expected {
				t.Errorf("GetConnectionString() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid minimal config",
			config: Config{
				Database: "test_db",
			},
			wantError: false,
		},
		{
			name:      "missing database",
			config:    Config{},
			wantError: true,
			errorMsg:  "MONGO_DATABASE is required",
		},
		{
			name: "AI enabled with OpenAI provider and API key",
			config: Config{
				Database:     "test_db",
				AIEnabled:    true,
				AIProvider:   "openai",
				OpenAIAPIKey: "test-key",
			},
			wantError: false,
		},
		{
			name: "AI enabled with OpenAI provider but no API key",
			config: Config{
				Database:   "test_db",
				AIEnabled:  true,
				AIProvider: "openai",
			},
			wantError: true,
			errorMsg:  "OPENAI_API_KEY is required when AI_PROVIDER is openai",
		},
		{
			name: "AI enabled with Gemini provider and API key",
			config: Config{
				Database:     "test_db",
				AIEnabled:    true,
				AIProvider:   "gemini",
				GeminiAPIKey: "test-key",
			},
			wantError: false,
		},
		{
			name: "AI enabled with invalid provider",
			config: Config{
				Database:   "test_db",
				AIEnabled:  true,
				AIProvider: "invalid",
			},
			wantError: true,
			errorMsg:  "invalid AI_PROVIDER: invalid (must be openai, gemini, or claude)",
		},
		{
			name: "Google Docs enabled with credentials path",
			config: Config{
				Database:              "test_db",
				GoogleDocsEnabled:     true,
				GoogleCredentialsPath: "/path/to/creds.json",
			},
			wantError: false,
		},
		{
			name: "Google Docs enabled without credentials",
			config: Config{
				Database:          "test_db",
				GoogleDocsEnabled: true,
			},
			wantError: true,
			errorMsg:  "either GOOGLE_CREDENTIALS_PATH or GOOGLE_CREDENTIALS_JSON is required when GOOGLE_DOCS_ENABLED is true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set up test environment variables
	testEnvVars := map[string]string{
		"MONGO_DATABASE":        "test_db",
		"MONGO_URL":             "mongodb://test:27017",
		"MONGO_USERNAME":        "testuser",
		"MONGO_PASSWORD":        "testpass",
		"AI_ENABLED":            "true",
		"AI_PROVIDER":           "openai",
		"OPENAI_API_KEY":        "test-openai-key",
		"MIGRATIONS_COLLECTION": "test_migrations",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up after test
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() failed: %v", err)
	}

	// Verify values were loaded correctly
	if cfg.Database != "test_db" {
		t.Errorf("Expected Database to be 'test_db', got '%s'", cfg.Database)
	}

	if cfg.MongoURL != "mongodb://test:27017" {
		t.Errorf("Expected MongoURL to be 'mongodb://test:27017', got '%s'", cfg.MongoURL)
	}

	if cfg.Username != "testuser" {
		t.Errorf("Expected Username to be 'testuser', got '%s'", cfg.Username)
	}

	if !cfg.AIEnabled {
		t.Error("Expected AIEnabled to be true")
	}

	if cfg.AIProvider != "openai" {
		t.Errorf("Expected AIProvider to be 'openai', got '%s'", cfg.AIProvider)
	}

	if cfg.MigrationsCollection != "test_migrations" {
		t.Errorf("Expected MigrationsCollection to be 'test_migrations', got '%s'", cfg.MigrationsCollection)
	}
}

func TestConfigDefaults(t *testing.T) {
	// Set only required environment variable
	os.Setenv("MONGO_DATABASE", "test_db")
	defer os.Unsetenv("MONGO_DATABASE")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() failed: %v", err)
	}

	// Check default values
	expectedDefaults := map[string]interface{}{
		"MongoURL":             "mongodb://localhost:27017",
		"MigrationsPath":       "./migrations",
		"MigrationsCollection": "schema_migrations",
		"AIProvider":           "openai",
		"AIEnabled":            false,
		"OpenAIModel":          "gpt-4o-mini",
		"GeminiModel":          "gemini-1.5-flash",
		"ClaudeModel":          "claude-3-5-sonnet-20241022",
		"MaxPoolSize":          10,
		"MinPoolSize":          1,
		"MaxIdleTime":          300,
		"Timeout":              60,
	}

	// Use reflection to check defaults (simplified version)
	if cfg.MongoURL != expectedDefaults["MongoURL"] {
		t.Errorf("Expected MongoURL default to be %s, got %s", expectedDefaults["MongoURL"], cfg.MongoURL)
	}

	if cfg.MigrationsCollection != expectedDefaults["MigrationsCollection"] {
		t.Errorf("Expected MigrationsCollection default to be %s, got %s", expectedDefaults["MigrationsCollection"], cfg.MigrationsCollection)
	}

	if cfg.AIProvider != expectedDefaults["AIProvider"] {
		t.Errorf("Expected AIProvider default to be %s, got %s", expectedDefaults["AIProvider"], cfg.AIProvider)
	}

	if cfg.MaxPoolSize != expectedDefaults["MaxPoolSize"] {
		t.Errorf("Expected MaxPoolSize default to be %d, got %d", expectedDefaults["MaxPoolSize"], cfg.MaxPoolSize)
	}
}
