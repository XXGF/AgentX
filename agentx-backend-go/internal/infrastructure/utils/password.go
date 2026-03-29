package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordUtils 密码工具类（对应 Java 的 PasswordUtils）

// EncodePassword 加密密码
func EncodePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MatchesPassword 验证密码是否匹配
func MatchesPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
