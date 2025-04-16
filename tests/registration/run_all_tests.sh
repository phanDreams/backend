#!/bin/bash
chmod +x tests/registration/*.sh

echo "Running registration tests..."
echo "----------------------------"

echo "\n1. Testing successful registration:"
./tests/registration/success.sh
echo "\n"

echo "2. Testing invalid input:"
./tests/registration/invalid_input.sh
echo "\n"

echo "3. Testing password validation:"
./tests/registration/password_validation.sh
echo "\n"

echo "4. Testing duplicate phone:"
./tests/registration/duplicate_phone.sh
echo "\n"
