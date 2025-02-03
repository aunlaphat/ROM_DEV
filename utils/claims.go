package utils

import (
	"errors"
	"fmt"
	"strconv"
)

// GetUserInfoFromClaims extracts userID and roleID from JWT claims.
func GetUserInfoFromClaims(claims map[string]interface{}) (userID string, roleID int, err error) {
	// ✅ Extract userID
	userIDVal, ok := claims["userID"]
	if !ok {
		return "", 0, errors.New("userID is missing in token claims")
	}
	userID, ok = userIDVal.(string)
	if !ok || userID == "" {
		return "", 0, errors.New("invalid userID in token claims")
	}

	// ✅ Extract and convert roleID
	roleID, err = GetRoleIDFromClaims(claims)
	if err != nil {
		return "", 0, err
	}

	return userID, roleID, nil
}

// GetUserIDFromClaims extracts only the userID from JWT claims.
func GetUserIDFromClaims(claims map[string]interface{}) (string, error) {
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid userID in token claims")
	}
	return userID, nil
}

// GetRoleIDFromClaims extracts and ensures roleID is an int.
func GetRoleIDFromClaims(claims map[string]interface{}) (int, error) {
	roleIDVal, ok := claims["roleID"]
	if !ok {
		fmt.Println("🔴 roleID is missing in claims:", claims) // Debug log
		return 0, fmt.Errorf("roleID is missing in token claims")
	}

	switch v := roleIDVal.(type) {
	case int:
		fmt.Println("🟢 roleID is an int:", v) // Debug log
		return v, nil
	case float64: // JSON numbers default to float64 in Go
		fmt.Println("🟡 roleID is float64, converting:", v) // Debug log
		return int(v), nil
	case string:
		convertedRoleID, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("🔴 Failed to convert roleID string to int:", v) // Debug log
			return 0, fmt.Errorf("invalid roleID format in token claims: %s", v)
		}
		fmt.Println("🟢 roleID converted from string:", convertedRoleID) // Debug log
		return convertedRoleID, nil
	default:
		fmt.Println("🔴 Invalid roleID type:", v) // Debug log
		return 0, fmt.Errorf("invalid roleID type in token claims")
	}
}
