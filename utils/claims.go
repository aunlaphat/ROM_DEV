package utils

import "fmt"

// Helper function ดึงข้อมูล userID และ roleID จาก claims
func GetUserInfoFromClaims(claims map[string]interface{}) (userID string, roleID string, err error) {
	// ดึง userID
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		return "", "", fmt.Errorf("invalid user information in token")
	}

	// ดึง roleID
	roleID, ok = claims["roleID"].(string)
	if !ok || roleID == "" {
		return "", "", fmt.Errorf("invalid role information in token")
	}

	return userID, roleID, nil
}

// Helper function ดึงข้อมูล userID จาก claims
func GetUserIDFromClaims(claims map[string]interface{}) (string, error) {
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid user information in token")
	}
	return userID, nil
}

// Helper function ดึงข้อมูล userID และ role จาก claims
func GetRoleIDFromClaims(claims map[string]interface{}) (string, error) {
	roleID, ok := claims["roleID"].(string)
	if !ok || roleID == "" {
		return "", fmt.Errorf("invalid role in token")
	}

	return roleID, nil
}
