package main

import (
	"context"

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
