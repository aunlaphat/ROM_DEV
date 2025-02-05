package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
)

func JWTMiddleware(tokenAuth *jwtauth.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, claims, err := jwtauth.FromContext(c.Request.Context())

		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "❌ Unauthorized", "data": nil})
			c.Abort()
			return
		}

		// ✅ เซ็ต claims เข้าไปใน Context
		c.Set("jwt_claims", claims)
		c.Next()
	}
}
