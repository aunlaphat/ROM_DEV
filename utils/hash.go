package utils

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

// HashPassword ใช้ SHA-256 สำหรับ Hash รหัสผ่าน
func HashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func IsValidPassword(password string) bool {
	matched, _ := regexp.MatchString(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`, password)
	return matched
}
