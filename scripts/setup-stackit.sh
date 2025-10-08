#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 STACKIT MongoDB Migration Tool Setup${NC}"
echo "=============================================="
echo

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -f ".env.stackit.example" ]]; then
    echo -e "${RED}Error: Please run this script from the mongo-migration-tool directory${NC}"
    exit 1
fi

# Copy STACKIT configuration template
if [[ -f ".env" ]]; then
    echo -e "${YELLOW}Warning: .env file already exists${NC}"
    read -p "Do you want to backup existing .env and create new one? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        mv .env .env.backup.$(date +%Y%m%d_%H%M%S)
        echo -e "${GREEN}Backed up existing .env file${NC}"
    else
        echo "Setup cancelled. Existing .env file preserved."
        exit 0
    fi
fi

cp .env.stackit.example .env
echo -e "${GREEN}✅ Created .env from STACKIT template${NC}"

# Interactive configuration
echo
echo -e "${BLUE}📝 MongoDB Configuration${NC}"
echo "Please provide your STACKIT MongoDB connection details:"
echo

read -p "MongoDB Cluster URL (e.g., cluster.stackit.cloud): " CLUSTER_URL
read -p "Database Name: " DATABASE_NAME
read -p "Username: " USERNAME
read -s -p "Password: " PASSWORD
echo

# Construct connection string
MONGO_URL="mongodb+srv://${USERNAME}:${PASSWORD}@${CLUSTER_URL}/${DATABASE_NAME}?retryWrites=true&w=majority&authSource=admin&maxPoolSize=10&minPoolSize=2&maxIdleTimeMS=600000&serverSelectionTimeoutMS=60000"

# Update .env file
sed -i.bak "s|MONGO_URL=.*|MONGO_URL=${MONGO_URL}|g" .env
sed -i.bak "s|MONGO_DATABASE=.*|MONGO_DATABASE=${DATABASE_NAME}|g" .env
rm -f .env.bak

echo
echo -e "${GREEN}✅ Configuration updated successfully!${NC}"
echo

# Test connection
echo -e "${BLUE}🔍 Testing MongoDB Connection${NC}"
echo "Building migration tool..."

if make build >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Build successful${NC}"
    
    echo "Testing connection to STACKIT MongoDB..."
    if ./build/mongo-migrate version >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Connection to STACKIT MongoDB successful!${NC}"
    else
        echo -e "${RED}❌ Connection failed. Please check your credentials and network connectivity.${NC}"
        echo
        echo "Common issues:"
        echo "- Incorrect username/password"
        echo "- Network connectivity to STACKIT cloud"
        echo "- Database name doesn't exist"
        echo "- Firewall blocking MongoDB connections"
        echo
        echo "You can test the connection manually with:"
        echo "./build/mongo-migrate status"
    fi
else
    echo -e "${RED}❌ Build failed. Please check for Go compilation errors.${NC}"
fi

echo
echo -e "${BLUE}📋 Next Steps${NC}"
echo "1. Verify your connection: ./build/mongo-migrate status"
echo "2. Create your first migration: ./build/mongo-migrate create 'initial setup'"
echo "3. Run migrations: ./build/mongo-migrate up"
echo
echo -e "${YELLOW}💡 Tips for STACKIT MongoDB:${NC}"
echo "- Use SSL/TLS (already configured)"
echo "- Connection timeouts are set to 60 seconds for cloud latency"
echo "- Connection pooling is optimized for cloud deployment"
echo "- Monitor your connection limits in STACKIT dashboard"
echo
echo -e "${GREEN}🎉 Setup complete! Your migration tool is ready for STACKIT MongoDB.${NC}"
