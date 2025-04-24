#!/bin/bash

chmod +x tests/registration/*.sh

echo -e "\nRunning registration tests..."
echo -e "----------------------------\n"

echo "1. Testing successful registration:"
./tests/registration/success.sh
echo -e "\n"

echo "2. Testing invalid input:"
./tests/registration/invalid_input.sh
echo -e "\n"

echo "3. Testing password validation:"
./tests/registration/password_validation.sh
echo -e "\n"

echo "4. Testing duplicate phone:"
./tests/registration/duplicate_phone.sh
echo -e "\n"

echo "5. Testing duplicate email:"
./tests/registration/duplicate_email.sh
echo -e "\n"
