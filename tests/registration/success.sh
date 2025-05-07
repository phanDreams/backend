#!/bin/bash
BASE_URL=${BASE_URL:-http://localhost:3000}
response=$(curl -s -w "\n%{http_code}" -X POST $BASE_URL/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Олена",
    "family_name": "Коваленко",
    "phone": "+380679876546",
    "email": "ol@example.com",
    "password": "StrongPass2024!",
    "password_confirmation": "StrongPass2024!"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "Status Code: $http_code"
echo "Response: $body"

if [ "$http_code" -eq 201 ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
fi
