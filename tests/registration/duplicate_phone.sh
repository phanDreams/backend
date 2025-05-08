#!/bin/bash
BASE_URL=${BASE_URL:-http://localhost:3000}
echo "First registration..."
response=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "First",
    "family_name": "User",
    "phone": "+38 (050) 123-45-67",
    "email": "first@example.com",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }'
)
echo -e "\nTesting duplicate phone..."
response=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/api/v1/specialists/register \
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
if [ "$http_code" -eq 409 ]; then
    echo "✅ Duplicate phone test passed"
else
  echo "❌ Duplicate-phone test failed: expected 409, got $http_code"
 echo "Response: $body"
  exit 1
fi