package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret")

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Пример ввода данных (обычно это email/пароль)
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверяем пользователя (например, сверяем пароль с хранилищем)
	if creds.Username != "admin" || creds.Password != "password" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := app.generateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (app *application) generateJWT(userID string) (string, error) {
	// Определяем срок действия токена
	expirationTime := time.Now().Add(24 * time.Hour)

	// Создаем токен
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	return token.SignedString(jwtSecret)
}

func (app *application) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			// http.Redirect(w, r, "/loginpage", http.StatusSeeOther)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[7:] // Убираем "Bearer "

		// Проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			// http.Redirect(w, r, "/loginpage", http.StatusSeeOther)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Токен валиден, продолжаем
		next.ServeHTTP(w, r)
	})
}
