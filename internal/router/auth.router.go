package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mikromolekula2002/Testovoe/internal/controllers"
)

type AuthRouter struct {
	authController *controllers.AuthController
}

func newAuthRouter(authController *controllers.AuthController) *AuthRouter {
	return &AuthRouter{authController}
}

func (ar *AuthRouter) SetAuthRoutes(ctx *gin.RouterGroup) {
	ctx.POST("/createTokens", ar.authController.CreateTokens)
	ctx.POST("/refreshTokens", ar.authController.RefreshTokens)
}
