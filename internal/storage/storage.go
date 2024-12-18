package storage

import (
	"context"

	"github.com/dsemenov12/shorturl/internal/models"
)

type Storage interface {
	Bootstrap(ctx context.Context) error
	Set(ctx context.Context, shortKey string, url string) (string, error)
	Get(ctx context.Context, shortKey string) (string, string, bool, error)
	GetUserURL(ctx context.Context) (result []models.ShortURLItem, err error)
	Delete(ctx context.Context, shortKey string) error
}