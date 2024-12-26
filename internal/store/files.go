package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"

	"github.com/Reensef/golang-react-boolib/internal/db"
	"github.com/gofrs/uuid"
)

const (
	FilesBucketName = "files"
)

type File struct {
	ID          int64       `json:"id" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	Size        int64       `json:"size" validate:"required"`
	Creator     FileCreator `json:"creator" validate:"required"`
	UUID        uuid.UUID   `json:"uuid"`
	Tag         string      `json:"tag"`
	TagID       *int        `json:"-"`
	ContentType string      `json:"content_type"`
	Downloads   int64       `json:"downloads"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type FileCreator struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type FilesStore struct {
	sqlDB  *sql.DB
	blobDB db.BlobDB
}

func (s *FilesStore) Create(ctx context.Context, file *File, data multipart.File) error {
	query := `
		INSERT INTO files (name, size, created_by_user_id)
		VALUES ($1, $2, $3)
		RETURNING files.id,
			(SELECT username FROM users WHERE id = files.created_by_user_id),
    		files.uuid,
    		files.created_at,
    		files.updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	err := s.sqlDB.QueryRowContext(
		ctx,
		query,
		file.Name,
		file.Size,
		file.Creator.ID,
	).Scan(
		&file.ID,
		&file.Creator.Username,
		&file.UUID,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if file.TagID != nil {
		query = `
			INSERT INTO file_to_tags (file_id, tag_id)
			VALUES ($1, $2)
			RETURNING (SELECT name FROM files_tags WHERE id = $2)
		`
		err = s.sqlDB.QueryRowContext(
			ctx,
			query,
			file.ID,
			*file.TagID,
		).Scan(
			&file.Tag,
		)
		if err != nil {
			return err
		}
	}

	err = s.blobDB.UploadFile(
		context.Background(),
		FilesBucketName,
		file.UUID.String(),
		data,
		file.Size,
		file.ContentType,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *FilesStore) GetAccessLinkByID(ctx context.Context, id int64) (*url.URL, error) {
	query := `
		SELECT f.uuid
		FROM files f
		WHERE f.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	row := s.sqlDB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return &url.URL{}, errors.Join(ErrDataNotFound, row.Err())
	}

	var uuid string
	err := row.Scan(
		&uuid,
	)

	if err != nil {
		return &url.URL{}, err
	}

	presignedURL, err := s.blobDB.GetAccessLink(context.Background(), FilesBucketName, uuid)
	if err != nil {
		return &url.URL{}, err
	}

	return presignedURL, nil
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
		f.id, f.name, f.uuid, COALESCE(files_tags.name, ''), f.size, f.created_by_user_id, 
		users.username, f.downloads, f.created_at, f.updated_at
		FROM files f
		JOIN users ON f.created_by_user_id = users.id
		LEFT JOIN file_to_tags ON file_to_tags.file_id = f.id
		LEFT JOIN files_tags ON file_to_tags.tag_id = files_tags.id
		WHERE f.is_deleted = false
	`
	args := []interface{}{}

	if tagID != "" {
		query += " AND files_tags.id = $1"
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

func (s *FilesStore) DeleteByID(ctx context.Context, id int64) error {
	query := `
		UPDATE files SET is_deleted = true, updated_at = CURRENT_TIMESTAMP
		WHERE files.id = $1
		RETURNING uuid
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	row := s.sqlDB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return errors.Join(ErrDataNotFound, row.Err())
	}

	var uuid string
	err := row.Scan(
		&uuid,
	)
	if err != nil {
		return err
	}

	err = s.blobDB.DeleteFile(ctx, FilesBucketName, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *FilesStore) IncrementDownloadCountByID(ctx context.Context, id int64) error {
	query := `
		UPDATE files SET downloads = downloads + 1, updated_at = CURRENT_TIMESTAMP
		WHERE files.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	row := s.sqlDB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return errors.Join(ErrDataNotFound, row.Err())
	}

	return nil
}
