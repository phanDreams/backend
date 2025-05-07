#!/bin/bash

echo "🚀 Deploying Go Backend App (Production Mode)"

# Navigate to project folder
cd ~/backend || { echo "❌ Failed to enter project directory"; exit 1; }

# Pull latest changes (if using git)
# git pull origin main

# Set environment variables for production
export ENV=production
export SERVER_ADDRESS=:443
export TLS_ENABLED=true
export TLS_CERT_FILE=/home/ubuntu/certs/server.crt
export TLS_KEY_FILE=/home/ubuntu/certs/server.key
export JWT_SECRET=yourproductionjwtsecret

# Build the Go app
echo "🛠 Building the app..."
go build -o pethelp ./cmd/pethelp || { echo "❌ Build failed"; exit 1; }

# Kill any running instance on port 443
echo "🛑 Killing old server (if any)..."
PID=$(sudo lsof -t -i:443)
if [ -n "$PID" ]; then
    sudo kill -9 "$PID"
    echo "✅ Old server killed"
fi

# Start the app with sudo (port 443 needs root)
echo "🚀 Starting server on port 443..."
sudo ./pethelp &

# Wait a bit for server to boot
sleep 5

# Confirm server started
sudo netstat -tuln | grep 443

echo "✅ Deployment Complete!"
