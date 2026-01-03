package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mangrover007/banking-backend/internals/repository"
	"github.com/Mangrover007/banking-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authController struct {
	s services.AuthService
}

func NewAuthController(service services.AuthService) AuthController {
	return &authController{
		s: service,
	}
}

type User struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type Login struct {
	Email       string `json:"email" binding:"required_without=phone_number,email"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required_without=email"` // add custom validator for phone number?
}

func (c *authController) Register(ctx *gin.Context) {
	var body User
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// hash user password
	password := body.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	password = string(hashedPassword)

	user := repository.RegisterUserParams{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		Address:     body.Address,
		Password:    password,
	}

	// delegate
	registeredUser, err := c.s.Register(ctx, user)
	if err != nil {
		if errors.Is(err, services.ErrPhoneIsRegistered) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, registeredUser)
}

func (c *authController) Login(ctx *gin.Context) {
	var body Login
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// hash pasword
	password := body.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// login in db
	sid, err := c.s.Login(ctx, body.PhoneNumber, body.Email, password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credentials",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// set cookie
	cookie := http.Cookie{
		Name:     "sid",
		Value:    sid.String(),
		Path:     "/trans",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	}
	ctx.SetCookieData(&cookie)
	ctx.JSON(http.StatusOK, nil)
}

func (c *authController) Logout(ctx *gin.Context) {
	sidstr, ok := ctx.Params.Get("sid")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "no session id",
		})
		return
	}

	sid, err := uuid.Parse(sidstr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// logout from server
	err = c.s.Logout(ctx, sid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// unset cookie
	cookie := http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	ctx.SetCookieData(&cookie)
	ctx.JSON(http.StatusOK, nil)
}
