package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Password  []byte `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UsersStore struct {
	sqlDB *sql.DB
}

func (s *UsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.password, u.role, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	data := &User{}

	row := s.sqlDB.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, errors.Join(ErrDataNotFound, row.Err())
	}

	err := row.Scan(
		&data.ID,
		&data.Email,
		&data.Username,
		&data.Password,
		&data.Role,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.password, u.role, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	data := &User{}

	row := s.sqlDB.QueryRowContext(ctx, query, email)
	if row.Err() != nil {
		return nil, errors.Join(ErrDataNotFound, row.Err())
	}

	err := row.Scan(
		&data.ID,
		&data.Email,
		&data.Username,
		&data.Password,
		&data.Role,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password, role)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryDBTimeout)
	defer cancel()

	err := s.sqlDB.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
