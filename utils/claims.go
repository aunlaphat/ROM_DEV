package utils

import (
	"errors"
	"fmt"
)

// GetUserInfoFromClaims ดึง userID และ roleID จาก claims
// - userID เป็น required field และต้องเป็น string ที่ไม่ว่างเปล่า
// - roleID เป็น optional field และต้องเป็น string ถ้ามี
func GetUserInfoFromClaims(claims map[string]interface{}) (userID, roleID string, err error) {
	// ดึง userID
	userIDVal, ok := claims["userID"]
	if !ok {
		return "", "", errors.New("userID is missing in token claims")
	}
	userID, ok = userIDVal.(string)
	if !ok || userID == "" {
		return "", "", errors.New("invalid userID in token claims")
	}

	// ดึง roleID (optional)
	if roleIDVal, ok := claims["roleID"]; ok {
		if roleID, ok = roleIDVal.(string); !ok {
			return "", "", errors.New("invalid roleID in token claims")
		}
	}

	return userID, roleID, nil
}

// GetUserIDFromClaims ดึง userID เท่านั้นจาก claims
func GetUserIDFromClaims(claims map[string]interface{}) (string, error) {
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid user information in token")
	}
	return userID, nil
}

// GetRoleIDFromClaims ดึง roleID เท่านั้นจาก claims
func GetRoleIDFromClaims(claims map[string]interface{}) (string, error) {
	roleID, ok := claims["roleID"].(string)
	if !ok || roleID == "" {
		return "", fmt.Errorf("invalid role in token")
	}

	return roleID, nil
}
