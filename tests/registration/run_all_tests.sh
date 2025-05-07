#!/bin/bash

chmod +x *.sh

echo -e "\nRunning registration tests..."
echo -e "----------------------------\n"

echo "1. Testing successful registration:"
./success.sh
echo -e "\n"

echo "2. Testing invalid input:"
./invalid_input.sh
echo -e "\n"

echo "3. Testing password validation:"
./password_validation.sh
echo -e "\n"

echo "4. Testing duplicate phone:"
./duplicate_phone.sh
echo -e "\n"

echo "5. Testing duplicate email:"
./duplicate_email.sh
echo -e "\n"
