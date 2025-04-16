#!/bin/bash

echo "Testing invalid email format..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "not-an-email",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')


if [ "$http_code" -eq 400 ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
fi

echo "Status Code: $http_code"
echo "Response: $body"

echo "\nTesting missing required field..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "Status Code: $http_code"
echo "Response: $body"
