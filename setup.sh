#!/bin/bash
# Copyright 2025 Christopher O'Connell
# All rights reserved

set -e

echo "MCL Setup Script"
echo "================"
echo ""

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24+ first."
    exit 1
fi

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

echo "✓ Prerequisites found"
echo ""

# Build the mcl binary
echo "Building mcl binary..."
go build -o mcl .
echo "✓ mcl binary built"
echo ""

# Build Docker image
echo "Building Docker image..."
docker build -t mcl:latest docker/
echo "✓ Docker image built"
echo ""

# Copy config if it doesn't exist
if [ ! -f "$HOME/.mcl.yml" ]; then
    echo "Creating default config at ~/.mcl.yml..."
    cp .mcl.yml.example "$HOME/.mcl.yml"
    echo "✓ Config file created"
    echo ""
    echo "📝 Please edit ~/.mcl.yml to add your custom domains and folders"
else
    echo "✓ Config file already exists at ~/.mcl.yml"
fi
echo ""

# Ask about PATH installation
echo "Would you like to install mcl to /usr/local/bin? (requires sudo)"
echo "This will make 'mcl' available system-wide."
read -p "Install to /usr/local/bin? [y/N]: " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    sudo cp mcl /usr/local/bin/
    echo "✓ mcl installed to /usr/local/bin"
else
    echo "To use mcl, either:"
    echo "  1. Add $(pwd) to your PATH:"
    echo "     export PATH=\"\$PATH:$(pwd)\""
    echo "  2. Or run it directly:"
    echo "     $(pwd)/mcl"
fi
echo ""

echo "✅ Setup complete!"
echo ""
echo "Quick Start:"
echo "  mcl new \"your first task\"     # Create a new container"
echo "  mcl list                       # List containers"
echo "  mcl connect <name>             # Connect to a container"
echo ""
echo "For more info, see README.md"