#!/usr/bin/env bash

# tests/login.sh
# Usage: BASE_URL=http://localhost:3000 ./tests/login.sh

BASE_URL=${BASE_URL:-http://localhost:3000}

response=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/api/v1/specialists/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teslenko@example.com",
    "password": "StrongPass2024!!!"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "Status Code: $http_code"
echo "Response: $body"

if [ "$http_code" -eq 200 ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    exit 1
fi
