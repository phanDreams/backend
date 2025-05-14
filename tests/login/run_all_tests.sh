#!/bin/bash

chmod +x *.sh

echo -e "\nRunning login tests..."
echo -e "----------------------------\n"

echo "1. Testing successful login:"
./success.sh
echo -e "\n"

echo "2. Testing invalid email:"
./invalid_email.sh
echo -e "\n"

echo "3. Testing invalid passwor:"
./invalid_password.sh
echo -e "\n"

