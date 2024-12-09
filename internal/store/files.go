package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Reensef/golang-react-boolib/internal/db"
	"github.com/gofrs/uuid"
)

type File struct {
	ID        int64       `json:"id" validate:"required"`
	Name      string      `json:"name" validate:"required"`
	Size      int64       `json:"size" validate:"required"`
	Creator   FileCreator `json:"sender"`
	UUID      uuid.UUID   `json:"uuid"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type FileCreator struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type FilesStore struct {
	sqlDB  *sql.DB
	blobDB *db.BlobDB
}

func (s *FilesStore) GetByID(ctx context.Context, id int64) (*File, error) {
	query := `
		SELECT f.id, f.name, f.uuid, f.size, f.created_by_user_id, u.username, f.created_at, f.updated_at
		FROM files f
		JOIN users u ON f.created_by_user_id = u.id
		WHERE f.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	data := &File{}

	row := s.sqlDB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, errors.Join(ErrDataNotFound, row.Err())
	}

	err := row.Scan(
		&data.ID,
		&data.Name,
		&data.UUID,
		&data.Size,
		&data.Creator.ID,
		&data.Creator.Username,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}
