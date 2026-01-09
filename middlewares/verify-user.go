package middlewares

import (
	"net/http"

	"github.com/Mangrover007/banking-backend/internals/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func VerifyUser(query *repository.Queries) gin.HandlerFunc {
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
