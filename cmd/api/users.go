package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Reensef/golang-react-boolib/internal/store"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserID = errors.New("invalid user ID")
)

type CreateUserPayload struct {
	Name     string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,max=50"`
	Password string `json:"password" validate:"required,max=50"`
}

// TODO Передавать на бекенд уже хешированный пароль
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	user := store.User{
		Username: payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := app.store.Users.Create(r.Context(), &user); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidUserID)
		return
	}

	comments, err := app.store.Users.GetByID(r.Context(), ID)

	if errors.Is(err, store.ErrDataNotFound) {
		app.resourceNotFoundResponse(w, r, err)
		return
	} else if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, comments); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
