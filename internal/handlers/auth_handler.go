package handlers

import (
	"errors"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"pethelp-backend/internal/domain/models"
	"pethelp-backend/internal/domain/service"
)

type RegistrationRequest struct {
	Name                 string `json:"name" binding:"required"`
	FamilyName          string `json:"family_name" binding:"required"`
	Phone               string `json:"phone" binding:"required" example:"+38 (XXX) XXX-XX-XX"`
	Email               string `json:"email" binding:"required,email" example:"user@example.com"`
	Password            string `json:"password" binding:"required,min=12"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

func isValidPassword(password string) error {
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

func (r *RegistrationRequest) ValidatePassword() error {
	validate := validator.New()

	// Validate basic struct tags (e.g. required, email, eqfield)
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Check password complexity
	if err := isValidPassword(r.Password); err != nil {
		return err
	}

	// Ensure the password is not the same as the email
	if r.Password == r.Email {
		return errors.New("password cannot be the same as email")
	}

	return nil
}

//handle register specialist endpoint
func RegisterSpecialistHandler(authService *service.AuthService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RegistrationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		if err := request.ValidatePassword(); err != nil {
			logger.Error("Password validation failed", zap.Error(err))
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		exists, err := service.CheckEmailExists(authService.DB, request.Email)
		if err != nil {
			logger.Error("Failed to check email existence", zap.Error(err))
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if exists {
			c.JSON(409, gin.H{"error": "Email already registered"})
			return
		}

		exists, err = service.CheckPhoneExists(authService.DB, request.Phone)
		if err != nil {
			logger.Error("Failed to check phone existence", zap.Error(err))
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if exists {
			c.JSON(409, gin.H{"error": "Phone number already registered"})
			return
		}

		newSpecialist := &models.Specialist{
			Name:       request.Name,
			FamilyName: request.FamilyName,
			Phone:      request.Phone,
			Email:      request.Email,
			Password:   request.Password,
		}

		err = authService.RegisterSpecialist(newSpecialist)
		if err != nil {
			logger.Error("Failed to register specialist", zap.Error(err))
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		token, err := authService.GenerateToken(newSpecialist)
		if err != nil {
			logger.Error("Failed to generate token", zap.Error(err))
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(201, gin.H{
			"message": "Specialist registered successfully",
			"id": newSpecialist.ID,
			"token":   token,
		})
	}
}
