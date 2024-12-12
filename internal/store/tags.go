package store

import (
	"context"
	"database/sql"
	"errors"
)

type Tag struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type TagsStore struct {
	sqlDB *sql.DB
}

func (s *TagsStore) GetAll(ctx context.Context) ([]*Tag, error) {
	query := `
		SELECT t.id, t.name
		FROM files_tags t
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	datas := []*Tag{}

	rows, err := s.sqlDB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Join(ErrDataNotFound, err)
	}
	defer rows.Close()

	for rows.Next() {
		data := Tag{}
		err := rows.Scan(
			&data.ID,
			&data.Name,
		)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &data)
	}

	return datas, nil
}
