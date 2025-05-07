package authinfrastructure

import (
	"errors"
	dom "pethelp-backend/internal/domain/auth"
	"unicode"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)


type BcryptHasher struct{}
func NewBcryptHasher() dom.PasswordHasher { return &BcryptHasher{} }

func (h *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

func (h *BcryptHasher) Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}


type Validator struct {
	v *validator.Validate
}

func NewPasswordValidator() dom.PasswordValidator{
	return &Validator{ v: validator.New()}
}

type RegistrationRequest struct {
	Name                 string `json:"name" binding:"required,min=2"`
	FamilyName          string `json:"family_name" binding:"required,min=2"`
	Phone string `json:"phone" binding:"required,regexp=^\\+[0-9]{1,3}[0-9\\- ()]{7,}$"`
	Email               string `json:"email" binding:"required,email"`
	Password            string `json:"password" binding:"required,min=12"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

func (v *Validator)  ValidatePassword (password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

