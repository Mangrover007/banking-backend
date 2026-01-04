package routers

import (
	"net/http"

	"github.com/Mangrover007/banking-backend/controllers"
	"github.com/Mangrover007/banking-backend/internals/repository"
	"github.com/Mangrover007/banking-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func verifyUser(query *repository.Queries) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		sidstr, err := ctx.Cookie("sid")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "login first dummy",
			})
			return
		}

		sid, err := uuid.Parse(sidstr)
		session, err := query.FindSession(ctx, sid)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "login first dummy",
			})
			return
		}

		user, err := query.FindUserByID(ctx, session.FkUserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, nil)
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	})
}

func BankRouter(master *gin.Engine, query *repository.Queries, conn *pgx.Conn) {
	var bankService = services.NewBank(conn)
	var bankController = controllers.NewBankController(bankService)

	router := master.Group("/bank")

	// verify user session id AND attach user struct to context
	router.Use(verifyUser(query))

	router.POST("/deposit", bankController.Deposit)
	router.POST("/withdraw", bankController.Withdraw)
	router.POST("/transfer", bankController.Transfer)
	router.POST("/open-savings", bankController.OpenSavingsAccount)
	router.POST("/open-checking", bankController.OpenCheckingAccount)
}
