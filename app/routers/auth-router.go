package routers

import (
	"net/http"

	"github.com/Mangrover007/banking-backend/app/controllers"
	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/Mangrover007/banking-backend/app/services"
	"github.com/gin-gonic/gin"
)

func AuthRouter(master *gin.Engine, query *repository.Queries) *gin.RouterGroup {

	var authService = services.NewAuthService(query)
	var authController = controllers.NewAuthController(authService)

	router := master.Group("/auth")
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "working"})
	})
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/logout", authController.Logout)

	return router
}
