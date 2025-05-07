package auth

type Credentialed interface {
	GetEmail() string
	SetPasswordHash(string)
	GetPasswordHash() string
  }
  