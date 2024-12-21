package store

import (
	"context"
	"database/sql"
	"errors"
)

type FileActionLog struct {
	ID         int64           `json:"id"`
	UserID     int64           `json:"user_id"`
	UserName   string          `json:"user_name"`
	FileID     int64           `json:"files_id"`
	FileName   string          `json:"file_name"`
	ActionID   FileActionLogID `json:"-"`
	ActionName string          `json:"action_name"`
	CreatedAt  string          `json:"created_at"`
}

type FilesActionsLogStore struct {
	sqlDB *sql.DB
}

type FileActionLogID int64

const (
	FileActionDownloaded = 1
	FileActionUploaded   = 2
	FileActionOpened     = 3
	FileActionDeleted    = 4
)

func (s *FilesActionsLogStore) Create(ctx context.Context, action *FileActionLog) error {
	query := `
		INSERT INTO files_actions_log (user_id, file_id, action_id)
		VALUES ($1, $2, $3)
		RETURNING id, (SELECT username FROM users WHERE id = $1), (SELECT name FROM files WHERE id = $2), (SELECT name FROM files_actions WHERE id = $3), created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	err := s.sqlDB.QueryRowContext(
		ctx,
		query,
		action.UserID,
		action.FileID,
		action.ActionID,
	).Scan(
		&action.ID,
		&action.UserName,
		&action.FileName,
		&action.ActionName,
		&action.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// TODO Добавить пагинацию
func (s *FilesActionsLogStore) GetAll(ctx context.Context) ([]*FileActionLog, error) {
	query := `
		SELECT 
			files_actions_log.id,
			files_actions_log.user_id,
			users.username,
			files_actions_log.file_id,
			files.name, 
			files_actions_log.action_id, 
			files_actions.name, 
			files_actions_log.created_at
		FROM files_actions_log
		JOIN users ON files_actions_log.user_id = users.id
		JOIN files ON files_actions_log.file_id = files.id
		JOIN files_actions ON files_actions_log.action_id = files_actions.id
		ORDER BY created_at DESC;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	datas := []*FileActionLog{}

	rows, err := s.sqlDB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Join(ErrDataNotFound, err)
	}
	defer rows.Close()

	for rows.Next() {
		data := FileActionLog{}
		err := rows.Scan(
			&data.ID,
			&data.UserID,
			&data.UserName,
			&data.FileID,
			&data.FileName,
			&data.ActionID,
			&data.ActionName,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &data)
	}

	return datas, nil
}
