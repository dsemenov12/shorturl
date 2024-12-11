package storage

import (
	"context"
	"database/sql"
)

type Storage interface {
	Bootstrap(ctx context.Context) error
	Set(ctx context.Context, shortKey string, url string) (string, error)
	Get(ctx context.Context, shortKey string) (string, string, bool, error)
	GetUserURL(ctx context.Context) (*sql.Rows, error)
	Delete(ctx context.Context, shortKey string) (result sql.Result, err error)
}