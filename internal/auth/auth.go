package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	secret := os.Getenv("SECRET_KEY")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return t
}

func VerifyToken(token string) bool {
	secret := os.Getenv("SECRET_KEY")
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false
	}
	return true
}
