#!/bin/bash
set -e

# Configuration
DOCKER_IMAGE="${DOCKER_IMAGE:-mongo-migration-tool}"
DOCKER_TAG="${BUILD_NUMBER:-latest}"
MIGRATION_TIMEOUT="${MIGRATION_TIMEOUT:-300}"  # 5 minutes default timeout
REQUIRE_SIGNED_IMAGES="${REQUIRE_SIGNED_IMAGES:-false}"
SIGNING_KEY_NAME="${SIGNING_KEY_NAME:-}"
KEY_VAULT_NAME="${KEY_VAULT_NAME:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] SUCCESS: $1${NC}"
}

# Function to check if migrations directory has changes
check_migration_changes() {
    log "Checking for migration changes..."
    
    # In a CI/CD environment, you might compare with the previous commit
    if [ -n "$CI" ]; then
        if git diff --quiet HEAD~1 HEAD -- migrations/; then
            log "No migration changes detected"
            return 1
        else
            log "Migration changes detected"
            return 0
        fi
    else
        # For local testing, always assume changes exist
        log "Running in local mode, assuming migrations need to be checked"
        return 0
    fi
}

build_migration_container() {
    log "Building migration container..."
    
    docker build -t "${DOCKER_IMAGE}:${DOCKER_TAG}" .
    
    if [ $? -eq 0 ]; then
        success "Migration container built successfully"
    else
        error "Failed to build migration container"
    fi
}

check_migration_status() {
    log "Checking migration status..."
    
    docker run --rm \
        --network="${DOCKER_NETWORK:-host}" \
        -e MONGO_URL="${MONGO_URL}" \
        -e MONGO_DATABASE="${MONGO_DATABASE}" \
        -e MONGO_USERNAME="${MONGO_USERNAME}" \
        -e MONGO_PASSWORD="${MONGO_PASSWORD}" \
        -e MIGRATIONS_COLLECTION="${MIGRATIONS_COLLECTION:-schema_migrations}" \
        "${DOCKER_IMAGE}:${DOCKER_TAG}" status
}

run_migrations() {
    log "Running migrations..."

    timeout "${MIGRATION_TIMEOUT}" docker run --rm \
        --network="${DOCKER_NETWORK:-host}" \
        -e MONGO_URL="${MONGO_URL}" \
        -e MONGO_DATABASE="${MONGO_DATABASE}" \
        -e MONGO_USERNAME="${MONGO_USERNAME}" \
        -e MONGO_PASSWORD="${MONGO_PASSWORD}" \
        -e MIGRATIONS_COLLECTION="${MIGRATIONS_COLLECTION:-schema_migrations}" \
        "${DOCKER_IMAGE}:${DOCKER_TAG}" up
    
    if [ $? -eq 0 ]; then
        success "Migrations completed successfully"
    else
        error "Migrations failed or timed out"
    fi
}

# Function to validate environment variables
validate_environment() {
    log "Validating environment variables..."
    
    if [ -z "$MONGO_URL" ]; then
        error "MONGO_URL environment variable is required"
    fi
    
    if [ -z "$MONGO_DATABASE" ]; then
        error "MONGO_DATABASE environment variable is required"
    fi
    
    log "Environment validation passed"
}


# Function to wait for MongoDB to be ready
wait_for_mongo() {
    log "Waiting for MongoDB to be ready..."

    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        log "Testing MongoDB connectivity, attempt $attempt/$max_attempts"

        if docker run --rm \
            --network="${DOCKER_NETWORK:-host}" \
            -e MONGO_URL="${MONGO_URL}" \
            -e MONGO_DATABASE="${MONGO_DATABASE}" \
            -e MONGO_USERNAME="${MONGO_USERNAME}" \
            -e MONGO_PASSWORD="${MONGO_PASSWORD}" \
            -e MONGO_SSL_ENABLED="${MONGO_SSL_ENABLED:-true}" \
            -e MONGO_TIMEOUT="${MONGO_TIMEOUT:-60}" \
            -e MIGRATIONS_COLLECTION="${MIGRATIONS_COLLECTION:-schema_migrations}" \
            "${DOCKER_IMAGE}:${DOCKER_TAG}" version >/dev/null 2>&1; then
            success "MongoDB is ready (cloud provider: $(echo $MONGO_URL | grep -o '[^@]*\.[^/]*' | tail -1))"
            return 0
        fi
        
        if [ $attempt -eq 1 ]; then
            log "Note: Testing connectivity to cloud MongoDB"
            log "This may take longer than local connections"
        fi
        
        sleep 3
        attempt=$((attempt + 1))
    done
    
    error "MongoDB failed to become ready within timeout. Please check:"
    echo "  - MongoDB URL: ${MONGO_URL}"
    echo "  - Network connectivity to STACKIT cloud"
    echo "  - Database credentials and permissions"
    echo "  - SSL/TLS configuration"
}

# Main execution
main() {
    log "Starting migration deployment process..."
    
    # Validate required environment variables
    validate_environment
    
    # Wait for MongoDB to be ready
    wait_for_mongo
    
    # Check if we need to run migrations
    if check_migration_changes || [ "$FORCE_MIGRATIONS" = "true" ]; then
        # Build the container (only if not using pre-built signed image)
        if [ "$REQUIRE_SIGNED_IMAGES" = "true" ]; then
            log "Using pre-built signed image: ${DOCKER_IMAGE}:${DOCKER_TAG}"
            # Verify the signed image
            verify_image_signature
        else
            build_migration_container
        fi
        
        # Show current status
        log "Current migration status:"
        check_migration_status || warn "Failed to get migration status (this might be expected for new databases)"
        
        # Run migrations
        run_migrations
        
        # Show final status
        log "Final migration status:"
        check_migration_status
        
        success "Migration deployment completed successfully"
    else
        log "No migration changes detected, skipping deployment"
    fi
}

# Parse command line arguments
case "${1:-auto}" in
    "auto")
        main
        ;;
    "force")
        export FORCE_MIGRATIONS=true
        main
        ;;
    "status")
        validate_environment
        wait_for_mongo
        if [ "$REQUIRE_SIGNED_IMAGES" = "true" ]; then
            verify_image_signature
        else
            build_migration_container
        fi
        check_migration_status
        ;;
    "build")
        build_migration_container
        ;;
    *)
        echo "Usage: $0 {auto|force|status|build}"
        echo "  auto   - Run migrations only if changes detected (default)"
        echo "  force  - Always run migrations regardless of changes"
        echo "  status - Show migration status only"
        echo "  build  - Build container only"
        exit 1
        ;;
esac
