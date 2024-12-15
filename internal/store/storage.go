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

type SortDirection int

const (
	AscendingOrder SortDirection = iota
	DescendingOrder
	NoOrder
)

type Storage struct {
	Files interface {
		GetByID(context.Context, int64) (*File, error)
		GetAll(ctx context.Context, sortBy string, sortDirection SortDirection, tagID string) ([]*File, error)
	}
	Tags interface {
		GetAll(context.Context) ([]*Tag, error)
	}
	Users interface {
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Create(context.Context, *User) error
	}
}

func NewStorage(sqlDB *sql.DB, blobDB *db.BlobDB) Storage {
	return Storage{
		Files: &FilesStore{sqlDB, blobDB},
		Tags:  &TagsStore{sqlDB},
		Users: &UsersStore{sqlDB},
	}
}
