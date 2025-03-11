package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
)

// ✅ JWT Middleware ตรวจสอบ Token และดึง `UserID` & `RoleID`
func JWTMiddleware(tokenAuth *jwtauth.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🔹 ดึง Token จาก Header หรือ Cookie
		tokenString, source := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
			c.Abort()
			return
		}

		// 🔹 ถอดรหัส Token
		claims, err := parseToken(c, tokenAuth, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token"})
			c.Abort()
			return
		}

		// 🔹 ดึงค่า `UserID` จาก Claims
		userID, err := getUserIDFromClaims(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 🔹 ดึงค่า `RoleID` จาก Claims
		roleID, err := getRoleIDFromClaims(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// ✅ เซ็ต `UserID` และ `RoleID` ใน Context
		c.Set("UserID", userID)
		c.Set("RoleID", roleID)

		// ✅ Debug Mode - แสดง Claims ทั้งหมด
		fmt.Printf("🔍 JWT Debug - UserID=%s, RoleID=%d, Claims=%v\n", userID, roleID, claims)

		// ✅ เซ็ตค่า Claims ตามแหล่งที่มา
		if source == "header" {
			c.Set("jwt_claims_header", claims)
		} else if source == "cookie" {
			c.Set("jwt_claims_cookie", claims)
		}
		c.Set("jwt_source", source)

		// 🔹 ส่งต่อไปยัง API
		c.Next()
	}
}

// ✅ ดึง Token จาก Header หรือ Cookie
func extractToken(c *gin.Context) (string, string) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(authHeader[7:]), "header"
	}

	token, err := c.Cookie("jwt")
	if err == nil && token != "" {
		return strings.TrimSpace(token), "cookie"
	}

	return "", ""
}

// ✅ ถอดรหัส Token และดึง Claims
func parseToken(c *gin.Context, tokenAuth *jwtauth.JWTAuth, tokenString string) (map[string]interface{}, error) {
	token, err := tokenAuth.Decode(tokenString)
	if err != nil {
		return nil, err
	}

	// ✅ Debug พิมพ์ค่า Claims ออกมา
	claims, err := token.AsMap(c.Request.Context())
	if err != nil {
		return nil, errors.New("unauthorized - invalid token claims format")
	}

	return claims, nil
}

// ✅ ดึงค่า `UserID` จาก Claims
func getUserIDFromClaims(claims map[string]interface{}) (string, error) {
	userID, exists := claims["userID"].(string)
	if !exists {
		return "", errors.New("unauthorized - missing UserID in token")
	}
	return userID, nil
}

// ✅ ดึงค่า `RoleID` จาก Claims
func getRoleIDFromClaims(claims map[string]interface{}) (int, error) {
	roleID, exists := claims["roleID"].(float64) // JSON Decode มาเป็น float64
	if !exists {
		return 0, errors.New("unauthorized - missing RoleID in token")
	}
	return int(roleID), nil
}
