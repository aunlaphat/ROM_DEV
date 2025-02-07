package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
)

func JWTAuthMiddleware(tokenAuth *jwtauth.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var source string

		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
			source = "header"
		}

		if tokenString == "" {
			token, err := c.Cookie("jwt")
			if err == nil {
				tokenString = token
				source = "cookie"
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
			c.Abort()
			return
		}

		claims, err := tokenAuth.Decode(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token"})
			c.Abort()
			return
		}

		if source == "header" {
			c.Set("jwt_claims_header", claims)
		} else if source == "cookie" {
			c.Set("jwt_claims_cookie", claims)
		}
		c.Set("jwt_source", source)

		c.Next()
	}
}
