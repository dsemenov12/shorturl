package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/storage/storage_main"
)

type ShortURLJSON struct {
	UUID 	    string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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
			UUID: strconv.Itoa(iter),
			ShortURL: key,
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

func Load(storage storage_main.Storage) error {
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