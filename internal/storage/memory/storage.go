package memory

import (
	"context"
	"sync"
	"database/sql"

	"github.com/dsemenov12/shorturl/internal/filestorage"
)

type StorageMemory struct {
	mx sync.RWMutex
    Data map[string]string
}

func NewStorage() *StorageMemory {
	StorageObj := StorageMemory{Data: make(map[string]string)}
    return &StorageObj
}

func (s *StorageMemory) Get(ctx context.Context, key string) (string, string, error) {
	s.mx.RLock()
    defer s.mx.RUnlock()
    return s.Data[key], key, nil
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

func (s *StorageMemory) GetUserURL(ctx context.Context) (rows *sql.Rows, err error) {
    return nil, nil
}

func (s *StorageMemory) Delete(ctx context.Context, shortKey string) (result sql.Result, err error) {
	delete(s.Data, shortKey)
	return nil, nil
}