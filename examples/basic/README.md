# Basic mongo-essential Library Example

This example demonstrates how to use mongo-essential as a Go library in your own projects.

## Setup

1. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` to match your MongoDB configuration:
   ```bash
   MONGO_URL=mongodb://localhost:27017
   MONGO_DATABASE=example_app
   ```

3. Make sure MongoDB is running locally:
   ```bash
   # Using Docker
   docker run --name mongo-example -p 27017:27017 -d mongo:7.0
   
   # Or use the make command from the root directory
   make db-up
   ```

4. Run the example:
   ```bash
   go mod init example
   go mod tidy
   go run main.go
   ```

## What this example demonstrates

- Loading configuration from environment variables and .env files
- Connecting to MongoDB with proper connection pooling and SSL/TLS support
- Creating and registering migrations
- Checking migration status
- Running migrations programmatically
- Implementing both Up and Down migration methods
- Error handling and logging

## Expected Output

```
Connected to MongoDB successfully
Migration Status:
  ❌ 20240101_001 - Create users collection with basic structure
  ❌ 20240101_002 - Add email and username indexes to users collection
Running pending migrations...
Creating users collection...
Adding indexes to users collection...
All migrations completed successfully!
```

## Integration in Your Project

To use mongo-essential in your own project:

1. Add the dependency:
   ```bash
   go get github.com/jocham/mongo-essential
   ```

2. Import the packages:
   ```go
   import (
       "github.com/jocham/mongo-essential/config"
       "github.com/jocham/mongo-essential/migration"
   )
   ```

3. Create your migrations following the pattern shown in `main.go`

4. Set up your configuration and run migrations as part of your application startup

## Advanced Usage

For more advanced features like AI analysis, cloud provider support, and Google Docs integration, see the main repository documentation and other examples.