package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test with environment variables
	_ = os.Setenv("MONGO_URL", "mongodb://testhost:27017")
	_ = os.Setenv("MONGO_DATABASE", "testdb")
	_ = os.Setenv("MIGRATIONS_COLLECTION", "test_migrations")
	_ = os.Setenv("OPENAI_API_KEY", "test-key")
	_ = os.Setenv("GEMINI_API_KEY", "gemini-key")
	_ = os.Setenv("GOOGLE_DRIVE_FOLDER_ID", "folder-123")

	defer func() {
		_ = os.Unsetenv("MONGO_URL")
		_ = os.Unsetenv("MONGO_DATABASE")
		_ = os.Unsetenv("MIGRATIONS_COLLECTION")
		_ = os.Unsetenv("OPENAI_API_KEY")
		_ = os.Unsetenv("GEMINI_API_KEY")
		_ = os.Unsetenv("GOOGLE_DRIVE_FOLDER_ID")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.MongoURL != "mongodb://testhost:27017" {
		t.Errorf("Expected MongoDB URI to be 'mongodb://testhost:27017', got '%s'", cfg.MongoURL)
	}

	if cfg.Database != "testdb" {
		t.Errorf("Expected MongoDB database to be 'testdb', got '%s'", cfg.Database)
	}

	if cfg.MigrationsCollection != "test_migrations" {
		t.Errorf("Expected migrations collection to be 'test_migrations', got '%s'", cfg.MigrationsCollection)
	}

	if cfg.OpenAIAPIKey != "test-key" {
		t.Errorf("Expected OpenAI API key to be 'test-key', got '%s'", cfg.OpenAIAPIKey)
	}

	if cfg.GeminiAPIKey != "gemini-key" {
		t.Errorf("Expected Gemini API key to be 'gemini-key', got '%s'", cfg.GeminiAPIKey)
	}

	if cfg.GoogleDriveFolderID != "folder-123" {
		t.Errorf("Expected Google Docs folder ID to be 'folder-123', got '%s'", cfg.GoogleDriveFolderID)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear all environment variables
	_ = os.Unsetenv("MONGO_URL")
	_ = os.Unsetenv("MONGO_DATABASE")
	_ = os.Unsetenv("MIGRATIONS_COLLECTION")
	_ = os.Unsetenv("OPENAI_API_KEY")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GOOGLE_DRIVE_FOLDER_ID")

	// Set required database field
	_ = os.Setenv("MONGO_DATABASE", "test")
	defer func() { _ = os.Unsetenv("MONGO_DATABASE") }()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.MongoURL != "mongodb://localhost:27017" {
		t.Errorf("Expected default MongoDB URI to be 'mongodb://localhost:27017', got '%s'", cfg.MongoURL)
	}

	if cfg.Database != "test" {
		t.Errorf("Expected MongoDB database to be 'test', got '%s'", cfg.Database)
	}

	if cfg.MigrationsCollection != "schema_migrations" {
		t.Errorf("Expected default migrations collection to be 'schema_migrations', got '%s'", cfg.MigrationsCollection)
	}

	if cfg.OpenAIAPIKey != "" {
		t.Errorf("Expected default OpenAI API key to be empty, got '%s'", cfg.OpenAIAPIKey)
	}

	if cfg.GeminiAPIKey != "" {
		t.Errorf("Expected default Gemini API key to be empty, got '%s'", cfg.GeminiAPIKey)
	}

	if cfg.GoogleDriveFolderID != "" {
		t.Errorf("Expected default Google Docs folder ID to be empty, got '%s'", cfg.GoogleDriveFolderID)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				MongoURL:             "mongodb://localhost:27017",
				Database:             "testdb",
				MigrationsCollection: "migrations",
			},
			expectError: false,
		},
		{
			name: "empty database",
			config: &Config{
				MongoURL:             "mongodb://localhost:27017",
				Database:             "",
				MigrationsCollection: "migrations",
			},
			expectError: true,
			errorMsg:    "MONGO_DATABASE is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected validation error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}
			}
		})
	}
}

// LoadFromFile tests removed - this API doesn't exist in current config implementation
