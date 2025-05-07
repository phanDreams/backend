package authinfrastructure

import (
	"net/http"

	appauth "pethelp-backend/internal/app/auth"
	dom "pethelp-backend/internal/domain/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RegistrationDTO interface {
	GetEmail() string
	GetPhone() string
	GetPassword() string
	GetPasswordConfirmation() string
	GetName() string
	GetFamilyName() string
}

//   type specialistReq struct {
// 		  Name                string `json:"name" binding:"required,min=2"`
// 		  FamilyName          string `json:"family_name" binding:"required,min=2"`
// 		  Phone 				string `json:"phone" binding:"required,regexp=^\\+[0-9]{1,3}[0-9\\- ()]{7,}$"`
// 		  Email               string `json:"email" binding:"required,email"`
// 		  Password            string `json:"password" binding:"required,min=12"`
// 		  PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
// }

func RegisterHandler[E dom.Registrable, DTO RegistrationDTO](authSvc   *appauth.AuthService,  
    validator dom.FieldsValidator,  
    repo      dom.Repository[E],    
    newDTO    func() DTO,           
    toEntity  func(DTO) E,         
    logger    *zap.Logger,) gin.HandlerFunc {
		return func(c *gin.Context) {
			// — Bind JSON into the right DTO
			req := newDTO()
			if err := c.ShouldBindJSON(req); err != nil {
				logger.Error("invalid payload", zap.Error(err))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			    return
			}
            // — Run shared field‐validator on it
			if err := validator.Validate(req); err != nil {
				logger.Warn("validation failed", zap.Error(err))
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			} 

			ctx := c.Request.Context()

			//email uniqueness
			if exists, err := repo.CheckEmailExists(ctx, req.GetEmail()); err != nil {
				logger.Error("email exists check failed", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			} else if exists {
				c.JSON(http.StatusConflict, gin.H{"error": dom.ErrAccountAlreadyExists.Error()})
				return
			}

			//phone uniqueness
			if exists, err := repo.CheckPhoneExists(ctx, req.GetPhone()); err != nil {
				logger.Error("phone exists check failed", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			} else if exists {
				c.JSON(http.StatusConflict, gin.H{"error": dom.ErrPhoneAlreadyInUse.Error()})
				return
			}

			// — Build domain entity
			entity := toEntity(req)
			// — Call the same Register for ANY Registrable
			if err := authSvc.Register(ctx, entity, entity, req.GetPassword()); err != nil {
				logger.Error("register failed", zap.Error(err))
				// you can switch on specific errors here…
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// — Success response
			c.JSON(http.StatusCreated, gin.H{
				"message": "Registration successful",
				"id":      entity.GetID(),
			})

		}

	}

