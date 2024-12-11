package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dsemenov12/shorturl/internal/auth"
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
			is_deleted boolean
		)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
    _, err = tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS short_key_idx ON storage (short_key)`)
	if err != nil {
		tx.Rollback()
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

func (s StorageDB) Get(ctx context.Context, shortKey string) (redirectLink string, shortKeyRes string, err error) {
    row := s.conn.QueryRowContext(ctx, "SELECT url, short_key FROM storage WHERE short_key=$1", shortKey)
	err = row.Scan(&redirectLink, &shortKeyRes)
	return
}

func (s StorageDB) GetUserURL(ctx context.Context) (rows *sql.Rows, err error) {
    return s.conn.QueryContext(ctx, "SELECT short_key, url FROM storage WHERE user_id=$1", ctx.Value(auth.UserIDKey))
}

func (s StorageDB) Delete(ctx context.Context, shortKey string) (result sql.Result, err error) {
	fmt.Println("UPDATE storage SET is_deleted=true WHERE short_key=$1")
	fmt.Println(shortKey)

    return s.conn.ExecContext(ctx, "UPDATE storage SET is_deleted=true WHERE short_key=$1", shortKey)
}