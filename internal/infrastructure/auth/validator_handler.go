package authinfrastructure

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	dom "pethelp-backend/internal/domain/auth"

	"github.com/go-playground/validator"
)

type fieldsValidatorImp struct {
	validate *validator.Validate
}

func NewFieldsValidator() dom.FieldsValidator{
	v := validator.New()
	   // register a no-arg "phone" validation
	   _ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
        // only digits/spaces/parens/dashes after the country code:
        re := regexp.MustCompile(`^\+[0-9]{1,3}[0-9\-\s()]{7,}$`)
        return re.MatchString(fl.Field().String())
    })
	return &fieldsValidatorImp{validate: v}
}

func (fv *fieldsValidatorImp) Validate(data interface{}) error {
	if err := fv.validate.Struct(data); err != nil {
		if inv, ok := err.(*validator.InvalidValidationError); ok {
			return inv
		}
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name", "FamilyName":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least 2 characters", err.Field()))
			case "Phone":
				errorMessages = append(errorMessages, "Phone must be in E.123 format (e.g., +38 (XXX) XXX-XX-XX)")
			case "Email":
				errorMessages = append(errorMessages, "Invalid email format")
			case "Password":
				errorMessages = append(errorMessages, "Password must be at least 12 characters")
			case "PasswordConfirmation":
				errorMessages = append(errorMessages, "Passwords do not match")
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()))
			}
		}
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return nil
}
