package auth

type PasswordValidator interface {
    ValidatePassword(plain string) error
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) error
  }
  