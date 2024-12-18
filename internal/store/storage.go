package store

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"net/url"
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

// TODO Добавить транзакции для совместных запросов sql и blob
type Storage struct {
	Files interface {
		Create(ctx context.Context, file *File, data multipart.File) error
		GetAccessLinkByID(context.Context, int64) (*url.URL, error)
		GetByID(context.Context, int64) (*File, error)
		GetAll(ctx context.Context, sortBy string, sortDirection SortDirection, tagID string) ([]*File, error)
		DeleteByID(context.Context, int64) error
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

func NewStorage(sqlDB *sql.DB, blobDB db.BlobDB) Storage {
	return Storage{
		Files: &FilesStore{sqlDB, blobDB},
		Tags:  &TagsStore{sqlDB},
		Users: &UsersStore{sqlDB},
	}
}
