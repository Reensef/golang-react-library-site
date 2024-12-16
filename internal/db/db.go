package db

import (
	"context"
	"database/sql"
	"io"
	"time"

	_ "github.com/lib/pq"
)

type BlobDB interface {
	UploadFile(ctx context.Context, bucketName, objectName string, file io.Reader, size int64, contentType string) error
	DownloadFile(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucketName, objectName string) error
}

func NewSql(
	addr string,
	maxOpenConns int,
	maxIdleConns int,
	maxIdleTime string,
) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewBlob(addr, accessKeyID, secretAccessKey string) (BlobDB, error) {
	var client BlobDB
	client, err := NewMinioBlobStore(addr, accessKeyID, secretAccessKey, false)
	if err != nil {
		return nil, err
	}

	return client, nil
}
