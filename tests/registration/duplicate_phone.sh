#!/bin/bash

echo "First registration..."
curl -s -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "First",
    "family_name": "User",
    "phone": "+38 (050) 123-45-67",
    "email": "first@example.com",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }'

echo "\nTesting duplicate phone..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Second",
    "family_name": "User",
    "phone": "+38 (050) 123-45-67",
    "email": "second@example.com",
    "password": "DifferentPass123!",
    "password_confirmation": "DifferentPass123!"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "Status Code: $http_code"
echo "Response: $body"
