package storage

import "context"

type Storage interface {
	Bootstrap(ctx context.Context) error
	Set(ctx context.Context, shortKey string, url string) (string, error)
	Get(ctx context.Context, shortKey string) (string, error)
}