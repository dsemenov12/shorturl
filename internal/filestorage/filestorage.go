package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/storage"
)

// ShortURLJSON представляет структуру данных для хранения сокращенного URL и его оригинала в JSON-формате.
type ShortURLJSON struct {
	UUID        string `json:"uuid"`         // Уникальный идентификатор для записи
	ShortURL    string `json:"short_url"`    // Сокращенный URL
	OriginalURL string `json:"original_url"` // Оригинальный URL
}

// Save сохраняет данные о сокращенных URL в файл в формате JSON.
func Save(storageData map[string]string) error {
	var iter = 1
	var data []byte

	file, err := os.OpenFile(config.FlagFileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range storageData {
		shortURLJSON := ShortURLJSON{
			UUID:        strconv.Itoa(iter),
			ShortURL:    key,
			OriginalURL: value,
		}

		data, err = json.Marshal(shortURLJSON)
		if err != nil {
			return err
		}
		data = append(data, '\n')

		iter++
	}

	_, err = file.Write(data)
	return err
}

// Load загружает данные из файла и сохраняет их в хранилище.
func Load(storage storage.Storage) error {
	var shortURLJSON *ShortURLJSON

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

		storage.Set(context.TODO(), shortURLJSON.ShortURL, shortURLJSON.OriginalURL)
	}

	return nil
}
