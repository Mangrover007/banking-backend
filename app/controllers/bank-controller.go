package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/Mangrover007/banking-backend/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BankController interface {
	Deposit(ctx *gin.Context)
	Withdraw(ctx *gin.Context)
	Transfer(ctx *gin.Context)
	OpenSavingsAccount(ctx *gin.Context)
	OpenCheckingAccount(ctx *gin.Context)
}

type controller struct {
	s services.BankService
}

func NewBankController(service services.BankService) BankController {
	return &controller{
		s: service,
	}
}

type AccountType string

var (
	AccountTypeSAVINGS  AccountType = "SAVINGS"
	AccountTypeCHECKING AccountType = "CHECKING"
)

var (
	AccTypeToRepoAccType = map[AccountType]repository.AccountType{
		AccountTypeSAVINGS:  repository.AccountTypeSAVINGS,
		AccountTypeCHECKING: repository.AccountTypeCHECKING,
	}
)

type reqDepositOrWithdraw struct {
	Amount int64       `json:"amount" binding:"required"`
	Type   AccountType `json:"type" binding:"required"`
}

type reqOpenAccount struct {
	Balance int64       `json:"balance" binding:"required"`
	Type    AccountType `json:"type" binding:"required"`
}

type reqTransfer struct {
	Amount    int64  `json:"amount" binding:"required"`
	Sender    string `json:"sender" binding:"required"`
	Recipient string `json:"recipient" binding:"required"`
}

func (a *AccountType) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}

	switch strings.ToUpper(str) {
	case "SAVINGS":
		*a = AccountTypeSAVINGS
	case "CHECKING":
		*a = AccountTypeCHECKING
	default:
		return errors.New("Not an account type")
	}

	return nil
}

func (c *controller) Deposit(ctx *gin.Context) {
	user := ctx.MustGet("user").(repository.User)
	var body reqDepositOrWithdraw
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	account, err := c.s.FindAccount(ctx, user.ID, AccTypeToRepoAccType[body.Type])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.s.Deposit(ctx, account, body.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *controller) Withdraw(ctx *gin.Context) {
	user := ctx.MustGet("user").(repository.User)
	var body reqDepositOrWithdraw
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	account, err := c.s.FindAccount(ctx, user.ID, AccTypeToRepoAccType[body.Type])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.s.Withdraw(ctx, account, body.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *controller) Transfer(ctx *gin.Context) {
	var body reqTransfer
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	senderID, err := uuid.Parse(body.Sender)
	sender, err := c.s.FindAccountByID(ctx, senderID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipientID, err := uuid.Parse(body.Recipient)
	recipient, err := c.s.FindAccountByID(ctx, recipientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.s.Transfer(ctx, sender, recipient, body.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *controller) OpenSavingsAccount(ctx *gin.Context) {
	user := ctx.MustGet("user").(repository.User)
	var body reqOpenAccount
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_type, _ := AccTypeToRepoAccType[body.Type]
	if _type != repository.AccountTypeSAVINGS {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "you are asshole",
		}) // asshole tried to send diff account type to diff endpoint
		return
	}

	account, err := c.s.OpenSavingsAccount(ctx, user, body.Balance, _type)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *controller) OpenCheckingAccount(ctx *gin.Context) {
	user := ctx.MustGet("user").(repository.User)
	var body reqOpenAccount
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_type, _ := AccTypeToRepoAccType[body.Type]
	if _type != repository.AccountTypeCHECKING {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "you are indeed an asshole",
		}) // asshole tried to send diff account type to diff endpoint
		return
	}

	account, err := c.s.OpenCheckingAccount(ctx, user, body.Balance, _type)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, account)
}
