package utils

import (
	"fmt"
	"strconv"
)

// GetUserInfoFromClaims extracts userID and roleID from JWT claims.
func GetUserInfoFromClaims(claims map[string]interface{}) (userID string, roleID int, err error) {
	// ✅ Extract userID
	userID, err = GetUserIDFromClaims(claims)
	if err != nil {
		return "", 0, fmt.Errorf("🔴 userID extraction failed: %w", err)
	}

	// ✅ Extract roleID
	roleID, err = GetRoleIDFromClaims(claims)
	if err != nil {
		return "", 0, fmt.Errorf("🔴 roleID extraction failed: %w", err)
	}

	return userID, roleID, nil
}

// GetUserIDFromClaims extracts only the userID from JWT claims.
func GetUserIDFromClaims(claims map[string]interface{}) (string, error) {
	// ✅ Ensure claims contain userID
	userIDVal, ok := claims["userID"]
	if !ok {
		return "", fmt.Errorf("🔴 userID is missing in token claims")
	}

	// ✅ Convert userID to string (handles different types)
	switch v := userIDVal.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("🔴 userID is empty in token claims")
		}
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return "", fmt.Errorf("🔴 invalid userID type in token claims")
	}
}

// GetRoleIDFromClaims extracts and ensures roleID is an int.
func GetRoleIDFromClaims(claims map[string]interface{}) (int, error) {
	roleIDVal, ok := claims["roleID"]
	if !ok {
		return 0, fmt.Errorf("🔴 roleID is missing in token claims")
	}

	switch v := roleIDVal.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		convertedRoleID, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("🔴 invalid roleID format in token claims: %s", v)
		}
		return convertedRoleID, nil
	default:
		return 0, fmt.Errorf("🔴 invalid roleID type in token claims")
	}
}
