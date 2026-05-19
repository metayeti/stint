#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Building stint for multiple platforms..."

# Create bin directory if it doesn't exist
mkdir -p bin

# --- 1. Windows (64-bit) ---
echo "  - Windows x64..."
GOOS=windows GOARCH=amd64 go build -o bin/stint-windows-amd64.exe ./cmd/stint

# --- 2. Linux (64-bit) ---
echo "  - Linux x64..."
GOOS=linux GOARCH=amd64 go build -o bin/stint-linux-amd64 ./cmd/stint

# --- 3. macOS (Intel) ---
echo "  - macOS x64 (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o bin/stint-darwin-amd64 ./cmd/stint

# --- 4. macOS (Apple Silicon M1/M2/M3) ---
echo "  - macOS arm64 (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o bin/stint-darwin-arm64 ./cmd/stint

echo "Done! Binaries are ready to go in the bin/ folder."