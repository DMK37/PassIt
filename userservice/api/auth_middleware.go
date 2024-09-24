package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing Authorization header"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		userID, err := parseToken(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createToken(userId, username string) (string, error) {

	claims := jwt.MapClaims{}

	claims["user_id"] = userId
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret")) // don't do this in production
}

func parseToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte("secret"), nil // don't do this in production
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}

	return userID, nil
}
