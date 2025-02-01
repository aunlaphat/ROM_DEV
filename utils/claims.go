package utils

import (
	"errors"
	"fmt"
)

// GetUserInfoFromClaims ดึง userID และ roleID จาก claims
// - userID เป็น required field และต้องเป็น string ที่ไม่ว่างเปล่า
// - roleID เป็น optional field และต้องเป็น string ถ้ามี
func GetUserInfoFromClaims(claims map[string]interface{}) (userID string, roleID int, err error) {
	// ดึง userID
	userIDVal, ok := claims["userID"]
	if !ok {
		return "", 0, errors.New("userID is missing in token claims")
	}
	userID, ok = userIDVal.(string)
	if !ok || userID == "" {
		return "", 0, errors.New("invalid userID in token claims")
	}

	// ดึง roleID (optional)
	roleIDVal, ok := claims["roleID"]
	if !ok {
		return userID, 0, nil // ถ้า roleID ไม่มี ก็ให้เป็นค่าเริ่มต้น 0
	}
	roleID, ok = roleIDVal.(int)
	if !ok {
		return "", 0, errors.New("invalid roleID in token claims")
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
func GetRoleIDFromClaims(claims map[string]interface{}) (int, error) {
	// ดึง roleID จาก claims
	roleIDVal, ok := claims["roleID"]
	if !ok {
		return 0, fmt.Errorf("roleID is missing in token claims")
	}

	// ตรวจสอบว่า roleID เป็น int หรือไม่
	roleID, ok := roleIDVal.(int)
	if !ok {
		// ถ้าไม่ใช่ int, ลองแปลงจาก string เป็น int
		if strRoleID, ok := roleIDVal.(string); ok {
			var convertedRoleID int
			// แปลง string เป็น int
			_, err := fmt.Sscanf(strRoleID, "%d", &convertedRoleID)
			if err != nil {
				return 0, fmt.Errorf("invalid roleID format in token claims")
			}
			return convertedRoleID, nil
		}
		// ถ้าไม่สามารถแปลงได้
		return 0, fmt.Errorf("invalid roleID format in token claims")
	}

	// ถ้าเป็น int ก็ให้ใช้เลย
	return roleID, nil
}
