package pg

import (
	"context"
	"database/sql"

	"github.com/dsemenov12/shorturl/internal/auth"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/config"
)

type StorageItem struct {
	UUID          string  `db:"user_id"`
	ShortURL      string  `db:"short_url"`
	OriginslURL   string  `db:"original_url"`
	DeletedFlag   bool    `db:"is_deleted"`
} 

type StorageDB struct {
    conn *sql.DB
}

func NewStorage(conn *sql.DB) *StorageDB {
    return &StorageDB{conn: conn}
}

func (s StorageDB) Bootstrap(ctx context.Context) error  {
    tx, err := s.conn.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
	defer tx.Rollback()

    _, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS storage(
			user_id varchar(36),
			short_key varchar(128),
			url text UNIQUE,
			is_deleted boolean DEFAULT false
		)
    `)
	if err != nil {
		return err
	}
    _, err = tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS short_key_idx ON storage (short_key)`)
	if err != nil {
		return err
	}

    return tx.Commit()
}

func (s StorageDB) Set(ctx context.Context, shortKey string, url string) (shortKeyResult string, err error) {
	_, err = s.conn.ExecContext(ctx, "INSERT INTO storage (short_key, url, user_id) VALUES ($1, $2, $3)", shortKey, url, ctx.Value(auth.UserIDKey))
	if err != nil {
		row := s.conn.QueryRowContext(ctx, "SELECT short_key FROM storage WHERE url=$1", url)
		row.Scan(&shortKeyResult)
	} else {
		shortKeyResult = shortKey
	}

	return shortKeyResult, err
}

func (s StorageDB) Get(ctx context.Context, shortKey string) (redirectLink string, shortKeyRes string, isDeleted bool, err error) {
    row := s.conn.QueryRowContext(ctx, "SELECT url, short_key, is_deleted FROM storage WHERE short_key=$1", shortKey)
	err = row.Scan(&redirectLink, &shortKeyRes, &isDeleted)
	return
}

func (s StorageDB) GetUserURL(ctx context.Context) (result []models.ShortURLItem, err error) {
	var shortKey string
	var originalURL string

    rows, err := s.conn.QueryContext(ctx, "SELECT short_key, url FROM storage WHERE user_id=$1", ctx.Value(auth.UserIDKey))
	if err != nil {
        return nil, err
    }

	for rows.Next() {
        err = rows.Scan(&shortKey, &originalURL)
        if err != nil {
           continue
        }

		result = append(result, models.ShortURLItem{
			OriginalURL: originalURL,
			ShortURL: config.FlagBaseAddr + "/" + shortKey,
		})
    }

	err = rows.Err()
    if err != nil {
        return nil, err
    }

	return result, nil
}

func (s StorageDB) Delete(ctx context.Context, shortKey string) (result sql.Result, err error) {
    return s.conn.ExecContext(ctx, "UPDATE storage SET is_deleted=true WHERE short_key=$1 AND user_id=$2", shortKey, ctx.Value(auth.UserIDKey))
}