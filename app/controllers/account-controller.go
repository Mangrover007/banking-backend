package controllers

import (
	"net/http"

	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/Mangrover007/banking-backend/app/services"
	"github.com/gin-gonic/gin"
)

type AccountController interface {
	GetUserAccounts(ctx *gin.Context)
}

type accountController struct {
	s services.AccountService
}

func NewAccountController(s services.AccountService) AccountController {
	return &accountController{
		s: s,
	}
}

func (c *accountController) GetUserAccounts(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	accounts, err := c.s.GetUserAccounts(ctx, user.(repository.User))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
