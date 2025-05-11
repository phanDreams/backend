package handlers

import (
	"fmt"
	"net/http"

	"pethelp-backend/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// OAuthHandlers OAuth2.0 services contains
type OAuthHandlers struct {
	OAuthService *service.OAuthUserService
}

// NewOAuthHandlers create new OAuthHandlers
func NewOAuthHandlers(service *service.OAuthUserService) *OAuthHandlers {
	return &OAuthHandlers{OAuthService: service}
}

// SignInWithProvider redirect to provider login page
func (h *OAuthHandlers) SignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler process provider callback and save user data
func (h *OAuthHandlers) ProviderCallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Printf("User received from provider: %+v\n", user)

	if h.OAuthService != nil {
		err = h.OAuthService.SetOrUpdateUser(c.Request.Context(), &user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save user: %w", err))
			return
		}
		fmt.Println("User saved to Redis successfully.")
	} else {
		fmt.Println("Warning: OAuthService is nil, user not saved to Redis.")
	}

	c.Redirect(http.StatusTemporaryRedirect, "/api/v1/oauth/success")
}

// Temporary success sign in handler
func (h *OAuthHandlers) Success(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`
      <div style="
          background-color: #fff;
          padding: 40px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          text-align: center;
      ">
          <h1 style="
              color: #333;
              margin-bottom: 20px;
          ">You have Successfully signed in!</h1>

          </div>
      </div>
  `)))
}
