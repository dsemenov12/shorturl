package memory

import (
	"context"
	"sync"

	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/models"
)

// StorageMemory представляет собой структуру для хранения данных в памяти.
type StorageMemory struct {
	mx   sync.RWMutex
	Data map[string]string
}

// NewStorage создает новый экземпляр StorageMemory с инициализацией пустой карты для хранения данных.
func NewStorage() *StorageMemory {
	StorageObj := StorageMemory{Data: make(map[string]string)}
	return &StorageObj
}

// Get извлекает из памяти значение для заданного ключа (сокращённого URL).
func (s *StorageMemory) Get(ctx context.Context, key string) (string, string, bool, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Data[key], key, false, nil
}

// Set сохраняет пару ключ-значение в память.
func (s *StorageMemory) Set(ctx context.Context, key string, value string) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Data[key] = value

	return value, nil
}

// Bootstrap загружает данные из внешнего хранилища в память, используя функционал файла.
func (s *StorageMemory) Bootstrap(ctx context.Context) error {
	filestorage.Load(s)
	return nil
}

// GetUserURL извлекает данные о всех сокращённых URL для текущего пользователя.
func (s *StorageMemory) GetUserURL(ctx context.Context) (result []models.ShortURLItem, err error) {
	return nil, nil
}

// Delete удаляет запись по ключу (сокращённому URL) из памяти.
func (s *StorageMemory) Delete(ctx context.Context, shortKey string) (err error) {
	delete(s.Data, shortKey)
	return nil
}
