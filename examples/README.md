# Migration Examples

This directory contains example migrations and a sample CLI application demonstrating how to use the mongo-essential migration library.

## Example Migrations

### 1. Add User Indexes (`20240101_001_add_user_indexes.go`)

Demonstrates how to:
- Create unique indexes
- Create compound indexes
- Set index options like background creation
- Handle index cleanup in rollback

**Operations:**
- Creates unique index on `email` field
- Creates descending index on `created_at` field
- Creates compound index on `status` and `created_at` fields

### 2. Transform User Data (`20240101_002_transform_user_data.go`)

Shows data transformation patterns:
- Iterating over collection documents
- Normalizing data (email to lowercase)
- Creating computed fields (`full_name` from `first_name` + `last_name`)
- Adding missing timestamps
- Conditional updates

### 3. Create Audit Collection (`20240101_003_create_audit_collection.go`)

Advanced collection operations:
- Creating collections with JSON Schema validation
- Setting up multiple indexes for different query patterns
- Using TTL indexes for automatic data expiration
- Schema validation with required fields and enums

## Sample CLI Application (`main.go`)

A complete example showing how to:
1. Load configuration using the config package
2. Connect to MongoDB
3. Create and configure a migration engine
4. Register multiple migrations
5. Implement up/down/status commands

### Usage

Make sure you have a MongoDB instance running and set your environment variables:

```bash
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DATABASE="example_db"
export MIGRATIONS_COLLECTION="schema_migrations"
```

Run the example commands:

```bash
# Show migration status
go run main.go status

# Run all pending migrations
go run main.go up

# Roll back the last migration
go run main.go down
```

### Expected Output

**Status command:**
```
Migration Status:
--------------------------------------------------------------------------------
Version              Applied    Applied At           Description
--------------------------------------------------------------------------------
20240101_001         ❌ No      Never                Add indexes to users collection for email and created_at fields
20240101_002         ❌ No      Never                Transform user data: normalize email case, add full_name field, and update timestamps
20240101_003         ❌ No      Never                Create audit collection with schema validation and indexes
```

**Up command:**
```
Running migrations up...
Running migration: 20240101_001 - Add indexes to users collection for email and created_at fields
✅ Completed migration: 20240101_001
Running migration: 20240101_002 - Transform user data: normalize email case, add full_name field, and update timestamps
✅ Completed migration: 20240101_002
Running migration: 20240101_003 - Create audit collection with schema validation and indexes
✅ Completed migration: 20240101_003
All migrations completed!
```

**Down command:**
```
Rolling back last migration...
Rolling back migration: 20240101_003 - Create audit collection with schema validation and indexes
✅ Rolled back migration: 20240101_003
```

## Integration with Your Project

To use these examples in your own project:

1. **Copy migration patterns**: Use the migration structs in `examplemigrations/` as templates
2. **Adapt the CLI**: Modify `main.go` to fit your application structure
3. **Configure properly**: Update import paths to match your module name
4. **Add error handling**: Enhance error handling as needed for production use
5. **Package structure**: The example migrations are in the `examplemigrations` package for proper Go module organization

## Testing the Examples

You can test the examples with a local MongoDB instance:

```bash
# Start MongoDB (if using Docker)
docker run -d -p 27017:27017 --name mongo-test mongo:latest

# Run the examples
cd examples
go run main.go status
go run main.go up
go run main.go status
go run main.go down
```

## Best Practices Demonstrated

1. **Idempotent operations**: Migrations handle cases where operations might be run multiple times
2. **Proper error handling**: Each migration returns errors appropriately
3. **Background index creation**: Uses background option to avoid blocking
4. **Schema validation**: Shows how to enforce data quality with JSON Schema
5. **Rollback strategies**: Each migration includes a proper down method
6. **Performance considerations**: Uses efficient query patterns and indexing