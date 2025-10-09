#!/bin/bash
set -e

echo "🚀 Setting up MCP migration testing tool..."

if [ ! -d "venv" ]; then
    echo "📦 Creating virtual environment..."
    python3 -m venv venv
fi

source venv/bin/activate

echo "🎉 Ready for testing"
