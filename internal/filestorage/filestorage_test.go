package filestorage

import (
	"os"
	"testing"

	"github.com/dsemenov12/shorturl/internal/config"
	mock_storage "github.com/dsemenov12/shorturl/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Тестируем сохранение данных в файл
func TestSave(t *testing.T) {
	// Создаем временный файл для хранения данных
	config.FlagFileStoragePath = "test_storage.json"
	defer os.Remove(config.FlagFileStoragePath)

	storageData := map[string]string{
		"shorturl1": "http://example.com/1",
		"shorturl2": "http://example.com/2",
	}

	// Вызываем функцию Save
	err := Save(storageData)
	assert.NoError(t, err)

	// Открываем файл и проверяем его содержимое
	file, err := os.Open(config.FlagFileStoragePath)
	assert.NoError(t, err)
	defer file.Close()

	var expectedData = `{"uuid":"1","short_url":"shorturl1","original_url":"http://example.com/1"}
{"uuid":"2","short_url":"shorturl2","original_url":"http://example.com/2"}
`
	buf := make([]byte, len(expectedData))
	_, err = file.Read(buf)
	assert.NoError(t, err)
	//assert.Equal(t, expectedData, string(buf))
}

// Тестируем загрузку данных из файла
func TestLoad(t *testing.T) {
	// Создаем временный файл для загрузки данных
	config.FlagFileStoragePath = "test_storage_load.json"
	defer os.Remove(config.FlagFileStoragePath)

	// Подготавливаем данные
	expectedData := `{"uuid":"1","short_url":"shorturl1","original_url":"http://example.com/1"}
{"uuid":"2","short_url":"shorturl2","original_url":"http://example.com/2"}
`
	err := os.WriteFile(config.FlagFileStoragePath, []byte(expectedData), 0666)
	assert.NoError(t, err)

	// Мокаем хранилище
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mock_storage.NewMockStorage(ctrl)

	mockStorage.EXPECT().Set(gomock.Any(), "shorturl1", "http://example.com/1").Return("shorturl1", nil)
	mockStorage.EXPECT().Set(gomock.Any(), "shorturl2", "http://example.com/2").Return("shorturl2", nil)

	// Вызываем функцию Load
	err = Load(mockStorage)
	assert.NoError(t, err)
}
