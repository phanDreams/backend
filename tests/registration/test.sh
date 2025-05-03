#!/bin/bash

echo "🧪 Running Local Tests for Go Backend"

# Navigate to project folder
cd ~/backend || { echo "❌ Failed to enter project directory"; exit 1; }

# Set environment variables for local testing
export ENV=local
export SERVER_ADDRESS=:3000
export TLS_ENABLED=false
export TLS_CERT_FILE=""
export TLS_KEY_FILE=""
export JWT_SECRET=testsecret

# Set BASE_URL for API tests
export BASE_URL=http://localhost:3000

# Build the Go app
echo "🛠 Building the app for testing..."
go build -o pethelp ./cmd/pethelp || { echo "❌ Build failed"; exit 1; }

# Kill old process on port 3000 if running
echo "🛑 Killing old local server (if any)..."
PID=$(lsof -t -i:3000)
if [ -n "$PID" ]; then
    kill -9 "$PID"
    echo "✅ Old server killed"
fi

# Start app in background
echo "🚀 Starting test server on port 3000..."
./pethelp &

# Capture background PID
SERVER_PID=$!

# Wait a bit for server to boot
sleep 3

# 🧪 Run Go unit tests
echo "🧪 Running Go unit tests..."
go test ./... || { echo "❌ Unit tests failed"; kill -9 "$SERVER_PID"; exit 1; }

# 🧪 Run API shell tests
echo "🧪 Running API shell tests..."
bash tests/registration/*.sh || { echo "❌ API tests failed"; kill -9 "$SERVER_PID"; exit 1; }

# Stop the server after tests
echo "🛑 Stopping test server..."
kill -9 "$SERVER_PID"

echo "✅ All Tests Completed Successfully!"
