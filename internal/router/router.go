package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/controllers"
)

type Router struct {
	Gin        *gin.Engine
	Config     *config.Config
	AuthRouter *AuthRouter
}

func NewRouter(config *config.Config, controller *controllers.Controller) *Router {
	ginRouter := gin.Default()

	ginRouter.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	return &Router{
		Gin:        ginRouter,
		Config:     config,
		AuthRouter: newAuthRouter(controller.AuthController),
	}
}

func (r *Router) SetRoutes() {
	auth := r.Gin.Group("/auth")

	r.AuthRouter.SetAuthRoutes(auth)

	if r.Config.EnvType != "prod" {
		//r.devRouter.setDevRoutes(auth)
		//r.Gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
