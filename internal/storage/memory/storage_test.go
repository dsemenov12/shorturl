package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageMemory_SetAndGet(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	// Test Set method
	key := "short1"
	value := "http://example.com"
	storedValue, err := storage.Set(ctx, key, value)
	assert.NoError(t, err, "Set should not return an error")
	assert.Equal(t, value, storedValue, "The stored value should match the expected value")

	// Test Get method
	gotValue, gotKey, _, err := storage.Get(ctx, key)
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, value, gotValue, "The retrieved value should match the expected value")
	assert.Equal(t, key, gotKey, "The retrieved key should match the expected key")
}

func TestStorageMemory_Delete(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	// Set a value
	key := "short1"
	value := "http://example.com"
	_, err := storage.Set(ctx, key, value)
	assert.NoError(t, err, "Set should not return an error")

	// Delete the key
	err = storage.Delete(ctx, key)
	assert.NoError(t, err, "Delete should not return an error")

	// Check if the key is deleted
	_, _, found, err := storage.Get(ctx, key)
	assert.NoError(t, err, "Get after delete should not return an error")
	assert.False(t, found, "Key should be deleted")
}

func TestStorageMemory_Bootstrap(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	// Test Bootstrap method (assuming filestorage.Load doesn't throw an error)
	err := storage.Bootstrap(ctx)
	assert.NoError(t, err, "Bootstrap should not return an error")
}

func TestStorageMemory_GetUserURL(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	// Test GetUserURL method (which currently returns nil)
	result, err := storage.GetUserURL(ctx)
	assert.NoError(t, err, "GetUserURL should not return an error")
	assert.Nil(t, result, "GetUserURL should return nil as there is no data yet")
}
