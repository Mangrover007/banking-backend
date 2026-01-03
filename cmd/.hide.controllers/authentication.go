package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mangrover007/banking-backend/cmd/entities"
	"github.com/Mangrover007/banking-backend/cmd/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type auth_controller struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) AuthController {
	return &auth_controller{
		service: service,
	}
}

func (c *auth_controller) Login(ctx *gin.Context) {
	var user entities.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	uid, err := c.service.Login(user)
	if err != nil {
		errNotFound := services.ErrUserNotFound
		if errors.Is(err, errNotFound) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		errInvalidCredentials := services.ErrInvalidCredentials
		if errors.Is(err, errInvalidCredentials) {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    uid,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
	}
	ctx.SetCookieData(&cookie)

	ctx.JSON(http.StatusOK, nil)
}

func (c *auth_controller) Register(ctx *gin.Context) {
	var user entities.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = c.service.Register(user)
	if err != nil {
		registerConflict := services.ErrRegisterConflict
		if errors.Is(err, registerConflict) {
			ctx.JSON(http.StatusConflict, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
