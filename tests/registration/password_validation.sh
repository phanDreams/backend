#!/bin/bash

echo "Testing too short password..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "test@example.com",
    "password": "Short1!",
    "password_confirmation": "Short1!"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "Status Code: $http_code"
echo "Response: $body"

echo "\nTesting password without special character..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "test@example.com",
    "password": "NoSpecialChar123",
    "password_confirmation": "NoSpecialChar123"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "Status Code: $http_code"
echo "Response: $body"
