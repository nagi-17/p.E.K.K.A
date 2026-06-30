package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func GenerateJWT(player_ID string) (string, error) {
	jwtConfig := config.LoadConfig()

	mapClaims := jwt.MapClaims{
		"player_id": player_ID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(jwtConfig.JWTSecret))
}
