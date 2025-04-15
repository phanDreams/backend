# Use the same email as above
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane",
    "family_name": "Smith",
    "phone": "+38 (050) 999-99-99",
    "email": "john.doe@example.com",
    "password": "DifferentPass123!",
    "password_confirmation": "DifferentPass123!"
  }'

# Expected: HTTP 409
# {
#   "error": "Email already registered"
# }