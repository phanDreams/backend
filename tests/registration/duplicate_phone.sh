# Use the same phone as first request
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice",
    "family_name": "Johnson",
    "phone": "+38 (050) 123-45-67",
    "email": "alice@example.com",
    "password": "DifferentPass123!",
    "password_confirmation": "DifferentPass123!"
  }'

# Expected: HTTP 409
# {
#   "error": "Phone number already registered"
# }