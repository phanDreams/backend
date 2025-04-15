#!/bin/bash

# Function to print colored output
print_test() {
    echo -e "\n\033[1;34m=== Running: $1 ===\033[0m"
}

print_result() {
    if [[ $1 -eq 0 ]]; then
        echo -e "\033[1;32m✓ Test passed\033[0m"
    else
        echo -e "\033[1;31m✗ Test failed\033[0m"
    fi
}

# Test successful registration
print_test "Successful Registration"
./success.sh
print_result $?

# Wait between tests
sleep 1

# Test duplicate email
print_test "Duplicate Email"
./duplicate_email.sh
print_result $?

sleep 1

# Test duplicate phone
print_test "Duplicate Phone"
./duplicate_phone.sh
print_result $?

sleep 1

# Test password validation
print_test "Password Validation"
./password_validation.sh
print_result $?

sleep 1

# Test invalid input
print_test "Invalid Input"
./invalid_input.sh
print_result $?

echo -e "\n\033[1;34m=== All tests completed ===\033[0m"