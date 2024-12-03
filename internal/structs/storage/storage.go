package storage

import (
	"sync"
	"os"
	"bufio"
	"encoding/json"

	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/config"
)

type Storage struct {
	mx sync.Mutex
    Data map[string]string
}

func NewStorage() *Storage {
	StorageObj := Storage{Data: make(map[string]string)}
    return &StorageObj
}

func (s *Storage) Get(key string) (string, error) {
	s.mx.Lock()
    defer s.mx.Unlock()
    return s.Data[key], nil
}

func (s *Storage) Set(key string, value string) {
	s.mx.Lock()
    defer s.mx.Unlock()
	s.Data[key] = value
}

func (s *Storage) Load() error {
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

		s.Set(shortURLJSON.ShortURL, shortURLJSON.OriginalURL)
	}

	return nil
}