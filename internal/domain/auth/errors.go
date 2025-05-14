package auth

import "errors"

var(
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenNotProvided = errors.New("token not provided")
	ErrTokenInvalid = errors.New("token invalid")
	ErrAuthFailed    = errors.New("authentication failed")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPhone = errors.New("invalid phone")
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidFamilyName = errors.New("invalid family name")
	ErrInvalidPasswordConfirmation = errors.New("invalid password confirmation")
	ErrInvalidPasswordLength = errors.New("invalid, password length must be between 8 and 12 characters")
	ErrInvalidPasswordSpecialCharacter = errors.New("invalid, password must contain special character")
	ErrInvalidPasswordNumber = errors.New("invalid, password must contain number")
	ErrInvalidPasswordUppercase = errors.New("invalid, password must contain uppercase")
	ErrInvalidPasswordLowercase = errors.New("invalid, password must contain lowercase")
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrPhoneAlreadyInUse = errors.New("phone already in use")
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrInvalidInput = errors.New("Invalid input")
)
