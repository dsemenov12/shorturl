package storage

import (
	"sync"
	"os"
	"bufio"
	"encoding/json"
	"context"

	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/config"
)

type StorageMemory struct {
	mx sync.Mutex
    Data map[string]string
}

func NewStorage() *StorageMemory {
	StorageObj := StorageMemory{Data: make(map[string]string)}
    return &StorageObj
}

func (s StorageMemory) Ping() error {
	s.mx.Lock()
    defer s.mx.Unlock()
    return nil
}

func (s *StorageMemory) Get(ctx context.Context, key string) (string, error) {
	s.mx.Lock()
    defer s.mx.Unlock()
    return s.Data[key], nil
}

func (s *StorageMemory) Set(ctx context.Context, key string, value string) (string, error) {
	s.mx.Lock()
    defer s.mx.Unlock()
	s.Data[key] = value

	return value, nil
}

func (s *StorageMemory) Bootstrap(ctx context.Context) error {
	var shortURLJSON *filestorage.ShortURLJSON

	file, err := os.OpenFile(config.FlagFileStoragePath, os.O_RDONLY, 0666)
    if err != nil {
        return err
    }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &shortURLJSON); err != nil {
			return err
		}

		s.Set(ctx, shortURLJSON.ShortURL, shortURLJSON.OriginalURL)
	}

	return nil
}