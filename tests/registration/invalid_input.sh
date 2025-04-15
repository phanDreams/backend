# Test invalid email format
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "not-an-email",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }'

# Expected: HTTP 400
# {
#   "error": "Invalid request body"
# }

# Test missing required field
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }'

# Expected: HTTP 400
# {
#   "error": "Invalid request body"
# }