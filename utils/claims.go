package utils

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// GetUserInfoFromClaims extracts userID and roleID from JWT claims.
func GetUserInfoFromClaims(claims map[string]interface{}, logger *zap.Logger) (userID string, roleID int, err error) {
	// ✅ Extract userID
	userID, err = GetUserIDFromClaims(claims)
	if err != nil {
		logger.Error("🔴 userID extraction failed", zap.Error(err))
		return "", 0, err
	}

	// ✅ Extract roleID
	roleID, err = GetRoleIDFromClaims(claims, logger)
	if err != nil {
		logger.Error("🔴 roleID extraction failed", zap.Error(err))
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
// GetRoleIDFromClaims extracts and ensures roleID is an int.
func GetRoleIDFromClaims(claims map[string]interface{}, logger *zap.Logger) (int, error) {
	roleIDVal, ok := claims["roleID"]
	if !ok {
		logger.Warn("🔴 roleID is missing in claims", zap.Any("claims", claims))
		return 0, fmt.Errorf("roleID is missing in token claims")
	}

	switch v := roleIDVal.(type) {
	case int:
		logger.Debug("🟢 roleID is an int", zap.Int("roleID", v))
		return v, nil
	case float64:
		logger.Debug("🟡 roleID is float64, converting", zap.Float64("roleID", v))
		return int(v), nil
	case string:
		convertedRoleID, err := strconv.Atoi(v)
		if err != nil {
			logger.Error("🔴 Failed to convert roleID string to int", zap.String("roleID", v))
			return 0, fmt.Errorf("invalid roleID format in token claims: %s", v)
		}
		logger.Debug("🟢 roleID converted from string", zap.Int("roleID", convertedRoleID))
		return convertedRoleID, nil
	default:
		logger.Error("🔴 Invalid roleID type", zap.Any("roleID", v))
		return 0, fmt.Errorf("invalid roleID type in token claims")
	}
}
