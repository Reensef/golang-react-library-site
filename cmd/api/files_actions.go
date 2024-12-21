package main

import (
	"context"
	"net/http"

	"github.com/Reensef/golang-react-boolib/internal/store"
)

func (app *application) logDeletedFile(ctx context.Context, userID int64, fileID int64) error {
	return app.store.FilesActionsLog.Create(ctx, &store.FileActionLog{
		UserID:   userID,
		FileID:   fileID,
		ActionID: store.FileActionDeleted,
	})
}

func (app *application) logOpenedFile(ctx context.Context, userID int64, fileID int64) error {
	return app.store.FilesActionsLog.Create(ctx, &store.FileActionLog{
		UserID:   userID,
		FileID:   fileID,
		ActionID: store.FileActionOpened,
	})
}

func (app *application) logUploadedFile(ctx context.Context, userID int64, fileID int64) error {
	return app.store.FilesActionsLog.Create(ctx, &store.FileActionLog{
		UserID:   userID,
		FileID:   fileID,
		ActionID: store.FileActionUploaded,
	})
}

func (app *application) logDownloadedFile(ctx context.Context, userID int64, fileID int64) error {
	return app.store.FilesActionsLog.Create(ctx, &store.FileActionLog{
		UserID:   userID,
		FileID:   fileID,
		ActionID: store.FileActionDownloaded,
	})
}

func (app *application) getFilesActionsLogHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(userInfoCtxKey).(UserCtxInfo)
	if !ok || userInfo.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	datas, err := app.store.FilesActionsLog.GetAll(r.Context())
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, datas); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
