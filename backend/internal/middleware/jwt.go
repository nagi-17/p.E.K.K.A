package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/config"
)

func VerifyJWT(next http.Handler) http.Handler {
	jwtConfig := config.LoadConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header in request", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid auth format: \"Bearer\" missing", http.StatusUnauthorized)
			return
		}
		if tokenStr == "" {
			http.Error(w, "token is missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, check_bool := token.Claims.(jwt.MapClaims)
		if check_bool == false {
			http.Error(w, "invalid payload", http.StatusUnauthorized)
			return
		}
		player_ID, check_bool := claims["player_id"].(string)
		if check_bool == false {
			http.Error(w, "player id is missing in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "player_id", player_ID)
		next.ServeHTTP(w, request.WithContext(ctx))
	})
}
