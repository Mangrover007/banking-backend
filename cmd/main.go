package main

import (
	"github.com/Mangrover007/banking-backend/cmd/controllers"
	"github.com/Mangrover007/banking-backend/cmd/entities"
	"github.com/Mangrover007/banking-backend/cmd/services"
	"github.com/gin-gonic/gin"
)

var (
	videoService services.VideoService       = services.New()
	controller   controllers.VideoController = controllers.New(videoService)
)

var (
	registeredDB = make(map[string]entities.User)
	activeUsers  = make(map[string]string)
)

var (
	authService services.AuthService       = services.NewAuthService(registeredDB, activeUsers)
	_           controllers.AuthController = controllers.NewAuthController(authService)
)

func main() {
	server := gin.New()

	server.Use(gin.Logger(), gin.Recovery())

	server.GET("/videos", controller.FindAll)
	server.POST("videos", controller.Save)

	server.Run(":8080")
}
