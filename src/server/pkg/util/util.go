package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateAccessToken(userID uint, role int) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"jti":  generateJTI(),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func GenerateRefreshToken(userID uint) (string, error) {
	secret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": generateJTI(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// ValidateAccessToken validates signature and expiry, returning all claims needed
// by the JWTAuth middleware to populate request context.
func ValidateAccessToken(tokenStr string) (userID uint, role int, jti string, expiry time.Time, err error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return 0, 0, "", time.Time{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, "", time.Time{}, fmt.Errorf("invalid claims")
	}

	sub, _ := claims["sub"].(float64)
	roleVal, _ := claims["role"].(float64)
	jtiVal, _ := claims["jti"].(string)
	expVal, _ := claims["exp"].(float64)
	return uint(sub), int(roleVal), jtiVal, time.Unix(int64(expVal), 0), nil
}
