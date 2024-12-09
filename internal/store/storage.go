package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Reensef/golang-react-boolib/internal/db"
)

var (
	ErrDataNotFound = errors.New("data not found")
	QueryDBTimeout  = time.Second * 5
)

type Storage struct {
	Files interface {
		GetByID(context.Context, int64) (*File, error)
	}
}

func NewStorage(sqlDB *sql.DB, blobDB *db.BlobDB) Storage {
	return Storage{
		Files: &FilesStore{sqlDB, blobDB},
	}
}
