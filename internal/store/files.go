package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Reensef/golang-react-boolib/internal/db"
	"github.com/gofrs/uuid"
)

type File struct {
	ID        int64       `json:"id" validate:"required"`
	Name      string      `json:"name" validate:"required"`
	Size      int64       `json:"size" validate:"required"`
	Creator   FileCreator `json:"creator" validate:"required"`
	UUID      uuid.UUID   `json:"uuid" validate:"required"`
	Tag       string      `json:"tag"`
	Downloads int64       `json:"downloads"`
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
		SELECT f.id, f.name, f.uuid, t.name, f.size, f.created_by_user_id, users.username, f.created_at, f.updated_at
		FROM files f
		JOIN users ON f.created_by_user_id = users.id
		JOIN file_to_tags ON file_to_tags.file_id = f.id
		JOIN files_tags ON file_to_tags.tag_id = files_tags.id
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
		&data.Tag,
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

func (s *FilesStore) GetAll(ctx context.Context, sortBy string, sortDirection SortDirection, tagID string) ([]*File, error) {
	query := `
		SELECT 
			f.id, f.name, f.uuid, files_tags.name, f.size, f.created_by_user_id, 
			users.username, f.downloads, f.created_at, f.updated_at
		FROM files f
		JOIN users ON f.created_by_user_id = users.id
		JOIN file_to_tags ON file_to_tags.file_id = f.id
		JOIN files_tags ON file_to_tags.tag_id = files_tags.id
	`
	args := []interface{}{}

	if tagID != "" {
		query += " WHERE files_tags.id = $1"
		args = append(args, tagID)
	}

	if sortBy != "" && sortDirection != NoOrder {
		dir := "ASC"
		if sortDirection == DescendingOrder {
			dir = "DESC"
		}

		validSortFields := map[string]bool{
			"name":       true,
			"size":       true,
			"created_at": true,
			"updated_at": true,
			"downloads":  true,
		}

		if validSortFields[sortBy] {
			query += " ORDER BY f." + sortBy + " " + dir
		} else {
			return nil, fmt.Errorf("invalid sortBy value: %s", sortBy)
		}
	}

	log.Println(query)
	log.Println(args)

	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	datas := []*File{}

	rows, err := s.sqlDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Join(ErrDataNotFound, err)
	}
	defer rows.Close()

	for rows.Next() {
		data := File{}
		err := rows.Scan(
			&data.ID,
			&data.Name,
			&data.UUID,
			&data.Tag,
			&data.Size,
			&data.Creator.ID,
			&data.Creator.Username,
			&data.Downloads,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &data)
	}

	return datas, nil
}
