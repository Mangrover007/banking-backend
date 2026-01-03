package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifySession(activeUsers map[string]string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("sid")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, nil)
			return
		}

		// call DB to verify is session is active
		// if not active, send fuck you
		_, ok := activeUsers[sessionID]
		if !ok {
			ctx.JSON(http.StatusForbidden, nil)
			return
		}

		// if active, attach user id to request?
		ctx.Next()
	}
}
