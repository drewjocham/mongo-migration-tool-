# mongo-essential as a Go Library

This document explains how to use mongo-essential as a Go library in your applications.

## Installation

Add mongo-essential to your Go project:

```bash
go get github.com/jocham/mongo-essential
```

## Basic Usage

### 1. Import the packages

```go
import (
    "github.com/jocham/mongo-essential/migration"
    "github.com/jocham/mongo-essential/config"
)
```

### 2. Load configuration

```go
// Load from .env files
cfg, err := config.Load(".env", ".env.local")
if err != nil {
    log.Fatal(err)
}

// Or load from environment only
cfg, err := config.LoadFromEnv()
if err != nil {
    log.Fatal(err)
}

// Validate configuration
if err := cfg.Validate(); err != nil {
    log.Fatal(err)
}
```

### 3. Connect to MongoDB

```go
import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

client, err := mongo.Connect(ctx, 
    options.Client().ApplyURI(cfg.GetConnectionString()))
if err != nil {
    log.Fatal(err)
}
db := client.Database(cfg.Database)
```

### 4. Create and run migrations

```go
// Create migration engine
engine := migration.NewEngine(db, cfg.MigrationsCollection)

// Register your migrations
engine.RegisterMany(
    &CreateUsersCollection{},
    &AddUserIndexes{},
    &CreateOrdersCollection{},
)

// Run pending migrations
if err := engine.Up(ctx, ""); err != nil {
    log.Fatal(err)
}
```

## Creating Migrations

Each migration must implement the `Migration` interface:

```go
type CreateUsersCollection struct{}

func (m *CreateUsersCollection) Version() string {
    return "20240101_001"  // Use YYYYMMDD_NNN format
}

func (m *CreateUsersCollection) Description() string {
    return "Create users collection with indexes"
}

func (m *CreateUsersCollection) Up(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("users")
    
    // Create indexes
    emailIndex := mongo.IndexModel{
        Keys:    bson.D{{"email", 1}},
        Options: options.Index().SetUnique(true),
    }
    
    _, err := collection.Indexes().CreateOne(ctx, emailIndex)
    return err
}

func (m *CreateUsersCollection) Down(ctx context.Context, db *mongo.Database) error {
    // Rollback logic
    collection := db.Collection("users")
    _, err := collection.Indexes().DropOne(ctx, "email_1")
    return err
}
```

## Migration Engine API

The migration engine provides these main methods:

### `NewEngine(db *mongo.Database, collection string) *Engine`
Creates a new migration engine.

### `Register(migration Migration)`
Registers a single migration.

### `RegisterMany(migrations ...Migration)`
Registers multiple migrations at once.

### `Up(ctx context.Context, target string) error`
Runs pending migrations forward. If `target` is empty, runs all pending migrations.

### `Down(ctx context.Context, target string) error`
Rolls back migrations to the target version.

### `GetStatus(ctx context.Context) ([]MigrationStatus, error)`
Returns the status of all migrations.

### `Force(ctx context.Context, version string) error`
Marks a migration as applied without running it.

## Configuration Options

### MongoDB Settings
- `MONGO_URL`: Connection URL (default: mongodb://localhost:27017)
- `MONGO_DATABASE`: Database name (required)
- `MONGO_USERNAME`: Authentication username
- `MONGO_PASSWORD`: Authentication password

### Migration Settings
- `MIGRATIONS_PATH`: Path to migrations directory (default: ./migrations)
- `MIGRATIONS_COLLECTION`: Collection name for tracking migrations (default: schema_migrations)

### SSL/TLS Settings
- `MONGO_SSL_ENABLED`: Enable SSL/TLS (default: false)
- `MONGO_SSL_INSECURE`: Skip certificate verification (default: false)

### Connection Pool Settings
- `MONGO_MAX_POOL_SIZE`: Maximum connection pool size (default: 10)
- `MONGO_MIN_POOL_SIZE`: Minimum connection pool size (default: 1)
- `MONGO_MAX_IDLE_TIME`: Maximum idle time in seconds (default: 300)
- `MONGO_TIMEOUT`: Connection timeout in seconds (default: 60)

### AI Analysis Settings (Optional)
- `AI_ENABLED`: Enable AI features (default: false)
- `AI_PROVIDER`: AI provider (openai, gemini, claude)
- `OPENAI_API_KEY`: OpenAI API key
- `GEMINI_API_KEY`: Google Gemini API key
- `CLAUDE_API_KEY`: Anthropic Claude API key

## Examples

See the `examples/` directory for complete working examples:

- `examples/basic/`: Basic migration setup and usage
- `examples/advanced/`: Advanced features and configurations

## Integration Patterns

### Application Startup

```go
func initDatabase() error {
    cfg, err := config.Load(".env")
    if err != nil {
        return err
    }
    
    client, err := mongo.Connect(context.Background(), 
        options.Client().ApplyURI(cfg.GetConnectionString()))
    if err != nil {
        return err
    }
    
    engine := migration.NewEngine(client.Database(cfg.Database), 
        cfg.MigrationsCollection)
    
    // Register all your migrations
    registerMigrations(engine)
    
    // Run migrations on startup
    return engine.Up(context.Background(), "")
}

func registerMigrations(engine *migration.Engine) {
    engine.RegisterMany(
        &migrations.CreateUsersCollection{},
        &migrations.AddUserIndexes{},
        &migrations.CreateOrdersCollection{},
        // ... add all your migrations here
    )
}
```

### Testing

```go
func TestMigrations(t *testing.T) {
    // Use testcontainers or mock database for testing
    cfg := &config.Config{
        MongoURL: "mongodb://localhost:27017",
        Database: "test_db",
        MigrationsCollection: "test_migrations",
    }
    
    client, err := mongo.Connect(context.Background(), 
        options.Client().ApplyURI(cfg.GetConnectionString()))
    require.NoError(t, err)
    
    engine := migration.NewEngine(client.Database(cfg.Database), 
        cfg.MigrationsCollection)
    
    // Test your migrations
}
```

## Error Handling

The library provides structured error handling:

```go
if err := engine.Up(ctx, ""); err != nil {
    // Handle migration errors
    log.Printf("Migration failed: %v", err)
    return err
}

// Check specific migration status
status, err := engine.GetStatus(ctx)
if err != nil {
    return err
}

for _, s := range status {
    if !s.Applied {
        log.Printf("Pending migration: %s - %s", s.Version, s.Description)
    }
}
```

## Best Practices

1. **Version Naming**: Use `YYYYMMDD_NNN` format for migration versions
2. **Idempotent Migrations**: Ensure migrations can be run multiple times safely
3. **Rollback Testing**: Always test your Down() methods
4. **Small Changes**: Keep migrations focused on single logical changes
5. **Documentation**: Add clear descriptions to your migrations
6. **Error Handling**: Handle errors appropriately in your migration code

For more detailed examples and advanced usage, see the complete examples in the `examples/` directory.