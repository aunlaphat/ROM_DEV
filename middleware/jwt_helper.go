package middleware

import "github.com/gin-gonic/gin"

// ✅ Helper Function ดึง JWT Claims ตาม Source (Header หรือ Cookie)
func GetJWTClaims(c *gin.Context) map[string]interface{} {
	source, exists := c.Get("jwt_source")
	if !exists {
		return nil
	}

	var claims interface{}
	if source == "header" {
		claims, _ = c.Get("jwt_claims_header")
	} else if source == "cookie" {
		claims, _ = c.Get("jwt_claims_cookie")
	}

	if claimsMap, ok := claims.(map[string]interface{}); ok {
		return claimsMap
	}
	return nil
}
