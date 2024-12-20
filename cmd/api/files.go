package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Reensef/golang-react-boolib/internal/store"
	"github.com/go-chi/chi/v5"
)

var (
	ErrFileUploadFailed = errors.New("failed upload file")
)

func (app *application) getFileHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidFileID)
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

func (app *application) getFilesHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortDirection := r.URL.Query().Get("sort_direction")

	var direction store.SortDirection = store.NoOrder
	if sortDirection == "asc" {
		direction = store.AscendingOrder
	} else if sortDirection == "desc" {
		direction = store.DescendingOrder
	}

	tagID := r.URL.Query().Get("tag_id")

	data, err := app.store.Files.GetAll(r.Context(), sortBy, direction, tagID)

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

func (app *application) downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(userInfoCtxKey).(UserCtxInfo)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidFileID)
		return
	}

	link, err := app.store.Files.GetAccessLinkByID(r.Context(), ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.logDownloadedFile(r.Context(), userInfo.ID, ID); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, link.String()); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) openFileHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(userInfoCtxKey).(UserCtxInfo)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidFileID)
		return
	}

	link, err := app.store.Files.GetAccessLinkByID(r.Context(), ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.logOpenedFile(r.Context(), userInfo.ID, ID); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, link.String()); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(userInfoCtxKey).(UserCtxInfo)
	if !ok || userInfo.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uploadFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fileTag := r.FormValue("tag")
	defer uploadFile.Close()

	file := store.File{
		Name: fileHeader.Filename,
		Size: fileHeader.Size,
		Creator: store.FileCreator{
			ID: userInfo.ID,
		},
		Tag:         fileTag,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	err = app.store.Files.Create(r.Context(), &file, uploadFile)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.logUploadedFile(r.Context(), userInfo.ID, file.ID); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, file); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(userInfoCtxKey).(UserCtxInfo)
	if !ok || userInfo.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidFileID)
		return
	}

	err = app.store.Files.DeleteByID(r.Context(), ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.logDeletedFile(r.Context(), userInfo.ID, ID); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
