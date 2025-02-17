package storage

import (
	"context"

	"github.com/dsemenov12/shorturl/internal/models"
)

// Storage определяет интерфейс для работы с хранилищем сокращенных URL-адресов.
type Storage interface {
	// Bootstrap инициализирует хранилище (например, создает таблицы в БД или загружает данные из файла).
	Bootstrap(ctx context.Context) error
	// Set сохраняет URL под заданным коротким ключом.
	// Если ключ уже существует, возвращает существующий ключ и ошибку.
	Set(ctx context.Context, shortKey string, url string) (string, error)
	// Get получает оригинальный URL по его короткому ключу.
	Get(ctx context.Context, shortKey string) (string, string, bool, error)
	// GetUserURL возвращает список всех URL, сохраненных пользователем.
	GetUserURL(ctx context.Context) (result []models.ShortURLItem, err error)
	// Delete помечает сокращенный URL как удаленный (soft delete).
	// Реальное удаление может происходить асинхронно.
	Delete(ctx context.Context, shortKey string) error
}
