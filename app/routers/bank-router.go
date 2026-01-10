package routers

import (
	"github.com/Mangrover007/banking-backend/app/controllers"
	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/Mangrover007/banking-backend/app/middlewares"
	"github.com/Mangrover007/banking-backend/app/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func BankRouter(master *gin.Engine, query *repository.Queries, conn *pgx.Conn) {
	var bankService = services.NewBank(conn)
	var bankController = controllers.NewBankController(bankService)

	router := master.Group("/bank")

	// verify user session id AND attach user struct to context
	router.Use(middlewares.VerifyUser(query))

	router.POST("/deposit", bankController.Deposit)
	router.POST("/withdraw", bankController.Withdraw)
	router.POST("/transfer", bankController.Transfer)
	router.POST("/open-savings", bankController.OpenSavingsAccount)
	router.POST("/open-checking", bankController.OpenCheckingAccount)
}
