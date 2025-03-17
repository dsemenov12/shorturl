package pg

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dsemenov12/shorturl/internal/auth"
	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestStorageDB_Bootstrap(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewStorage(db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS storage`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`CREATE UNIQUE INDEX IF NOT EXISTS short_key_idx`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = storage.Bootstrap(ctx)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorageDB_Set(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewStorage(db)
	ctx := context.WithValue(context.Background(), auth.UserIDKey, "test-user")

	shortKey := "short123"
	originalURL := "https://example.com"

	mock.ExpectExec("INSERT INTO storage").
		WithArgs(shortKey, originalURL, "test-user").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.Set(ctx, shortKey, originalURL)
	assert.NoError(t, err)
	assert.Equal(t, shortKey, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorageDB_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewStorage(db)
	ctx := context.Background()

	shortKey := "short123"
	originalURL := "https://example.com"

	mock.ExpectQuery("SELECT url, short_key, is_deleted FROM storage").
		WithArgs(shortKey).
		WillReturnRows(sqlmock.NewRows([]string{"url", "short_key", "is_deleted"}).
			AddRow(originalURL, shortKey, false))

	url, key, deleted, err := storage.Get(ctx, shortKey)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)
	assert.Equal(t, shortKey, key)
	assert.False(t, deleted)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorageDB_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewStorage(db)
	ctx := context.WithValue(context.Background(), auth.UserIDKey, "test-user")

	shortKey := "short123"

	mock.ExpectExec("UPDATE storage SET is_deleted=true").
		WithArgs(shortKey, "test-user").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.Delete(ctx, shortKey)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorageDB_GetUserURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewStorage(db)
	ctx := context.WithValue(context.Background(), auth.UserIDKey, "test-user")

	shortKey1 := "short123"
	originalURL1 := "https://example.com/1"
	shortKey2 := "short456"
	originalURL2 := "https://example.com/2"

	mock.ExpectQuery(`SELECT short_key, url FROM storage WHERE user_id=\$1`).
		WithArgs("test-user").
		WillReturnRows(sqlmock.NewRows([]string{"short_key", "url"}).
			AddRow(shortKey1, originalURL1).
			AddRow(shortKey2, originalURL2))

	result, err := storage.GetUserURL(ctx)
	assert.NoError(t, err)

	assert.Len(t, result, 2)
	assert.Equal(t, result[0].ShortURL, config.FlagBaseAddr+"/"+shortKey1)
	assert.Equal(t, result[0].OriginalURL, originalURL1)
	assert.Equal(t, result[1].ShortURL, config.FlagBaseAddr+"/"+shortKey2)
	assert.Equal(t, result[1].OriginalURL, originalURL2)

	assert.NoError(t, mock.ExpectationsWereMet())
}
