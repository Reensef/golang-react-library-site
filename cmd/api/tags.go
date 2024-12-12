package main

import (
	"errors"
	"net/http"

	"github.com/Reensef/golang-react-boolib/internal/store"
)

func (app *application) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.store.Tags.GetAll(r.Context())

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
