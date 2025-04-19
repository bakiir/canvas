package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("key") // Секретный ключ, желательно из .env!

func GenerateJWT(teacherID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"teacher_id": teacherID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(), // Жизнь токена = 24 часа
	})

	return token.SignedString(jwtKey)
}
