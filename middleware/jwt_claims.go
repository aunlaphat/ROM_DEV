package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

// JWTAuthMiddleware validates JWT and extracts claims
func JWTAuthMiddleware(tokenAuth *jwtauth.JWTAuth, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate and verify JWT token
		token, err := jwtauth.VerifyRequest(tokenAuth, c.Request, jwtauth.TokenFromHeader, jwtauth.TokenFromQuery)
		if err != nil {
			logger.Error("❌ Invalid JWT Token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Ensure token is present
		if token == nil {
			logger.Error("🚫 Missing JWT Token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		// Extract claims from token
		claims, _ := token.AsMap(c.Request.Context())

		// Set JWT claims in the context for further use
		c.Set("jwt_claims", claims)

		/* // Extract RoleID for Role-Based Access Control
		if roleID, ok := claims["roleID"].(float64); ok {
			c.Set("roleID", int(roleID)) // Convert float64 to int (JWT encodes numbers as float64)
		} else {
			logger.Warn("⚠️ roleID not found in token claims")
		} */

		c.Next()
	}
}
