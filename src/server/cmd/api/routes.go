package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/app"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/middleware"
)

func SetupRoutes(a *app.App) {
	a.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	docs := a.Router.Group("/docs")
	{
		docs.GET("", a.DocsController.SwaggerUI)
		docs.GET("/", a.DocsController.SwaggerUI)
		docs.GET("/openapi.yaml", a.DocsController.OpenAPISpec)
		docs.GET("/openapi.yml", a.DocsController.OpenAPISpec)
	}

	v1 := a.Router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", a.AuthController.SignUp)
			auth.POST("/sign-in", a.AuthController.SignIn)
			auth.POST("/sign-out",
				middleware.Auth(a.JWT, a.Denylist),
				a.AuthController.SignOut,
			)
		}

		profile := v1.Group("/profile", middleware.Auth(a.JWT, a.Denylist))
		{
			profile.GET("/:id", a.ProfileController.Get)
			profile.PATCH("", a.ProfileController.Update)
		}
	}
}
