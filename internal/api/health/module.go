package health

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const (
	healthRoutePath = "/health"
)

var Module = fx.Module("health",
	fx.Invoke(registerRoutes))

func registerRoutes(route *gin.Engine) {
	healthGroup := route.Group(healthRoutePath)
	{
		healthGroup.GET("", Check())
	}
}
