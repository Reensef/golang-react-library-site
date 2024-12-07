package store

import (
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
}

func NewStorage(sqlDB *sql.DB, blobDB *db.BlobDB) Storage {
	return Storage{}
}
