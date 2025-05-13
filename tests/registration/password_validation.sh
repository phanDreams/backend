#!/bin/bash
BASE_URL=${BASE_URL:-http://localhost:3000}
echo "Testing too short password..."
response=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/api/v1/specialists/register \
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
response=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/api/v1/specialists/register \
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
if [ "$http_code" -eq 400 ]; then
    echo "✅ Password validation test passed"
else
  echo "❌ Duplicate-phone test failed: expected 500, got $http_code"
 echo "Response: $body"
  exit 1
fi