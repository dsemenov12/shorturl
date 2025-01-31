package memory

import (
	"context"
	"sync"

	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/models"
)

type StorageMemory struct {
	mx sync.RWMutex
    Data map[string]string
}

func NewStorage() *StorageMemory {
	StorageObj := StorageMemory{Data: make(map[string]string)}
    return &StorageObj
}

func (s *StorageMemory) Get(ctx context.Context, key string) (string, string, bool, error) {
	s.mx.RLock()
    defer s.mx.RUnlock()
    return s.Data[key], key, false, nil
}

func (s *StorageMemory) Set(ctx context.Context, key string, value string) (string, error) {
	s.mx.Lock()
    defer s.mx.Unlock()
	s.Data[key] = value

	return value, nil
}

func (s *StorageMemory) Bootstrap(ctx context.Context) error {
	filestorage.Load(s)
	return nil
}

func (s *StorageMemory) GetUserURL(ctx context.Context) (result []models.ShortURLItem, err error) {
    return nil, nil
}

func (s *StorageMemory) Delete(ctx context.Context, shortKey string) (err error) {
	delete(s.Data, shortKey)
	return nil
}