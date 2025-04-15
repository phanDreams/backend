curl -X POST http://localhost:3000/api/v1/specialists/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "family_name": "Doe",
    "phone": "+38 (050) 123-45-67",
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "password_confirmation": "SecurePass123!"
  }'

# Expected: HTTP 201
# {
#   "message": "Specialist registered successfully",
#   "id": 1
# }