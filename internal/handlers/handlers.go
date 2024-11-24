package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/structs/storage"
	"github.com/go-chi/chi/v5"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("pgx", config.FlagDatabaseDSN)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
    }
    defer db.Close()

	if err := db.Ping(); err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
    }

	res.WriteHeader(http.StatusOK)
}

func ShortenPost(res http.ResponseWriter, req *http.Request) {
	var inputDataValue models.InputData

	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &inputDataValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	storage.StorageObj.Set(shortKey, inputDataValue.URL)

	var result = models.ResultJSON{
		Result: shortURL,
	}

	resp, err := json.MarshalIndent(result, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	if config.FlagDatabaseDSN != "" {
		db, err := sql.Open("pgx", config.FlagDatabaseDSN)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()
		
		_, err = db.ExecContext(req.Context(), "CREATE TABLE IF NOT EXISTS storage(short_key TEXT, url TEXT)")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		db.ExecContext(req.Context(), "INSERT INTO storage (short_key, url) VALUES ($1, $2)", shortKey, inputDataValue.URL)
	} else {
		filestorage.Save(storage.StorageObj.Data)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func PostURL(res http.ResponseWriter, req *http.Request) {
	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if config.FlagDatabaseDSN != "" {
		db, err := sql.Open("pgx", config.FlagDatabaseDSN)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()

		_, err = db.ExecContext(req.Context(), "CREATE TABLE IF NOT EXISTS storage(short_key TEXT, url TEXT)")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		db.ExecContext(req.Context(), "INSERT INTO storage (short_key, url) VALUES ($1, $2)", shortKey, string(body))
	} else {
		storage.StorageObj.Set(shortKey, string(body))
		filestorage.Save(storage.StorageObj.Data)
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")

	var redirectLink string
	var err error

	if config.FlagDatabaseDSN != "" {
		db, err := sql.Open("pgx", config.FlagDatabaseDSN)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()

		_, err = db.ExecContext(req.Context(), "CREATE TABLE IF NOT EXISTS storage(short_key TEXT, url TEXT)")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	
		row := db.QueryRowContext(req.Context(), "SELECT url FROM storage WHERE short_key=$1", shortKey)
		
		err = row.Scan(&redirectLink)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
	} else {
		redirectLink, err = storage.StorageObj.Get(shortKey)
		if err != nil {
			http.Error(res, "redirect not found", 404)
		}
	}

	http.Redirect(res, req, redirectLink, http.StatusTemporaryRedirect)
}
