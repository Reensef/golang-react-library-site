package main

import (
	"encoding/json"
	"net/http"

	"github.com/Reensef/golang-react-boolib/internal/auth"
)

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
	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (app *application) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[7:] // Убираем "Bearer "

		token, err := auth.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
