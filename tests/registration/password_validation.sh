# Test too short password
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "test@example.com",
    "password": "Short1!",
    "password_confirmation": "Short1!"
  }'

# Expected: HTTP 400
# {
#   "error": "Password must be at least 12 characters long and contain uppercase, lowercase, number, and special characters"
# }

# Test password without special character
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "test@example.com",
    "password": "NoSpecialChar123",
    "password_confirmation": "NoSpecialChar123"
  }'

# Expected: HTTP 400
# {
#   "error": "Password must be at least 12 characters long and contain uppercase, lowercase, number, and special characters"
# }

# Test password mismatch
curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "family_name": "User",
    "phone": "+38 (050) 111-11-11",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "password_confirmation": "DifferentPass123!"
  }'

# Expected: HTTP 400
# {
#   "error": "Password confirmation does not match"
# }