package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Reensef/golang-react-boolib/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Println(creds)

	user, err := app.store.Users.GetByEmail(r.Context(), creds.Email)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(creds.Password))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userInfo := UserCtxInfo{
			ID:   int64(claims["userID"].(float64)),
			Role: claims["role"].(string),
		}

		ctx := context.WithValue(r.Context(), userInfoCtxKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
