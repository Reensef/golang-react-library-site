package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Reensef/golang-react-boolib/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getFileHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidPostID)
		return
	}

	data, err := app.store.Files.GetByID(r.Context(), ID)

	if errors.Is(err, store.ErrDataNotFound) {
		app.resourceNotFoundResponse(w, r, err)
		return
	} else if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, data); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
