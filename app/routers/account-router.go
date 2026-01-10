package routers

import (
	"github.com/Mangrover007/banking-backend/app/controllers"
	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/Mangrover007/banking-backend/app/middlewares"
	"github.com/Mangrover007/banking-backend/app/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func AccountRouter(master *gin.Engine, conn *pgx.Conn, query *repository.Queries) {
	var service = services.NewAccountService(conn, query)
	var controller = controllers.NewAccountController(service)

	router := master.Group("/accounts")
	router.Use(middlewares.VerifyUser(query))
	router.GET("/all", controller.GetUserAccounts)
}
