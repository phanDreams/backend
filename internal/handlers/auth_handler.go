package handlers

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"pethelp-backend/internal/domain/models"
	"pethelp-backend/internal/domain/service"
)

type RegistrationRequest struct {
	Name                 string `json:"name" binding:"required,min=2"`
	FamilyName          string `json:"family_name" binding:"required,min=2"`
	Phone               string `json:"phone" binding:"required" example:"+38 (012) 345-67-89"`
	Email               string `json:"email" binding:"required,email"`
	Password            string `json:"password" binding:"required,min=12"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
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

func (r *RegistrationRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
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
			}
		}
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return nil
}

// RegisterSpecialistHandler handles register specialist endpoint
func RegisterSpecialistHandler(authService *service.AuthService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the raw body first
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			c.JSON(400, gin.H{"error": "Failed to read request"})
			return
		}
		// Restore the body for subsequent reading
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		var request RegistrationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Error("Failed to bind JSON", 
				zap.Error(err),
				zap.String("raw_body", string(bodyBytes)))
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := request.Validate(); err != nil {
			logger.Error("Validation failed", zap.Error(err))
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		exists, err := authService.CheckEmailExists(request.Email)
		if err != nil {
			logger.Error("Failed to check email existence", zap.Error(err))
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if exists {
			c.JSON(409, gin.H{"error": "Email already registered"})
			return
		}

		exists, err = authService.CheckPhoneExists(request.Phone)
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
			Is_banned:  false,
			Is_deleted: false,
			Is_active:  true,
			Is_verified: false,
		}

		err = authService.RegisterSpecialist(newSpecialist)
		if err != nil {
			logger.Error("Failed to register specialist", 
				zap.Error(err),
				zap.Any("specialist", newSpecialist)) // Add specialist data logging
			c.JSON(500, gin.H{"error": err.Error()}) // Return actual error
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
			"id":     newSpecialist.ID,
			"token":  token,
		})
	}
}
