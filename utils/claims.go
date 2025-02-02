package utils

import (
	"errors"
	"fmt"
	"strconv"
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
		// ลองแปลงจาก float64 หรือ string เป็น int
		switch v := roleIDVal.(type) {
		case float64:
			// แปลงจาก float64 เป็น int
			roleID = int(v)
		case string:
			// แปลงจาก string เป็น int
			roleID, err = strconv.Atoi(v)
			if err != nil {
				return "", 0, fmt.Errorf("invalid roleID format in token claims: %w", err)
			}
		default:
			return "", 0, errors.New("invalid roleID format in token claims")
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

// GetRoleIDFromClaims ดึง roleID จาก JWT Claims และแปลงให้ถูกต้อง
func GetRoleIDFromClaims(claims map[string]interface{}) (int, error) {
	roleIDVal, ok := claims["roleID"]
	if !ok {
		return 0, fmt.Errorf("roleID is missing in token claims")
	}

	// ✅ Debug Log ตรวจสอบค่าก่อนแปลง
	fmt.Printf("🔍 Debug: roleIDVal=%v (Type: %T)\n", roleIDVal, roleIDVal)

	switch v := roleIDVal.(type) {
	case int:
		return v, nil
	case float64: // ✅ JSON อาจเก็บ roleID เป็น float64 -> แปลงเป็น int
		return int(v), nil
	case string: // ✅ ถ้า roleID เป็น string -> แปลงเป็น int
		roleID, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid roleID format in token claims: %w", err)
		}
		return roleID, nil
	default:
		return 0, fmt.Errorf("invalid roleID format in token claims (type: %T)", roleIDVal)
	}
}
