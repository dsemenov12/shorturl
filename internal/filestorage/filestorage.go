package filestorage

import (
	"os"
	"encoding/json"
	"strconv"
	"bufio"
	"fmt"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/structs/storage"
)

type ShortUrlJSON struct {
	Uuid 	    string `json:"uuid"`
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
		shortUrlJSON := ShortUrlJSON{
			Uuid: strconv.Itoa(iter),
			ShortURL: key,
			OriginalURL: value,
		}

		data, err = json.Marshal(shortUrlJSON)
		if err != nil {
			return err
		}
		data = append(data, '\n')

		iter++
	}

	_, err = file.Write(data)
    return err
}

func Load() error {
	var shortUrlJSON *ShortUrlJSON

	file, err := os.OpenFile(config.FlagFileStoragePath, os.O_RDONLY, 0666)
    if err != nil {
        return err
    }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &shortUrlJSON); err != nil {
			return err
		}

		storage.StorageObj.Set(shortUrlJSON.ShortURL, shortUrlJSON.OriginalURL)
	}

	fmt.Println(storage.StorageObj.Data)

	return nil
}