package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTUtils JWT 工具类（对应 Java 的 org.xhy.infrastructure.utils.JwtUtils）
type JWTUtils struct {
	secret     string
	expiration time.Duration
}

// NewJWTUtils 创建 JWT 工具实例
func NewJWTUtils(secret string, expirationSeconds int64) *JWTUtils {
	return &JWTUtils{
		secret:     secret,
		expiration: time.Duration(expirationSeconds) * time.Second,
	}
}

// Claims 自定义 JWT Claims
type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token（对应 Java 的 JwtUtils.generateToken）
func (j *JWTUtils) GenerateToken(userID string) (string, error) {
	if userID == "" {
		zap.L().Error("生成JWT Token失败: 用户ID为空")
		return "", errors.New("用户ID不能为空")
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		zap.L().Error("生成JWT Token失败", zap.Error(err))
		return "", err
	}

	zap.L().Debug("成功生成JWT Token",
		zap.String("userId", userID),
		zap.Time("expiresAt", now.Add(j.expiration)),
	)
	return tokenString, nil
}

// GetUserIDFromToken 从 Token 中获取用户ID（对应 Java 的 JwtUtils.getUserIdFromToken）
func (j *JWTUtils) GetUserIDFromToken(tokenString string) (string, error) {
	if tokenString == "" {
		zap.L().Warn("获取用户ID失败: Token为空")
		return "", errors.New("Token为空")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		zap.L().Warn("解析Token失败", zap.Error(err))
		return "", err
	}

	if !token.Valid {
		zap.L().Warn("Token无效")
		return "", errors.New("Token无效")
	}

	zap.L().Debug("成功解析Token", zap.String("userId", claims.Subject))
	return claims.Subject, nil
}

// ValidateToken 验证 Token 是否有效（对应 Java 的 JwtUtils.validateToken）
func (j *JWTUtils) ValidateToken(tokenString string) bool {
	if tokenString == "" {
		zap.L().Debug("Token验证失败: Token为空")
		return false
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		zap.L().Warn("Token验证失败", zap.Error(err))
		return false
	}

	if !token.Valid {
		return false
	}

	zap.L().Debug("Token验证成功",
		zap.String("userId", claims.Subject),
	)
	return true
}
