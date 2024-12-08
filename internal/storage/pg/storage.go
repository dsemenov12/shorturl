package pg

import (
	"context"
	"database/sql"
)

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
			short_key varchar(128),
			url TEXT UNIQUE,
			user_id varchar(36)
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
	_, err = s.conn.ExecContext(ctx, "INSERT INTO storage (short_key, url, user_id) VALUES ($1, $2, $3)", shortKey, url, ctx.Value("user_id"))
	if err != nil {
		row := s.conn.QueryRowContext(ctx, "SELECT short_key FROM storage WHERE url=$1 AND user_id=$2", url, ctx.Value("user_id"))
		row.Scan(&shortKeyResult)
	} else {
		shortKeyResult = shortKey
	}

	return shortKeyResult, err
}

func (s StorageDB) Get(ctx context.Context, shortKey string) (redirectLink string, err error) {
    row := s.conn.QueryRowContext(ctx, "SELECT short_key, url FROM storage WHERE short_key=$1", shortKey)
	err = row.Scan(&redirectLink)
	return
}

func (s StorageDB) GetUserURL(ctx context.Context) (rows *sql.Rows, err error) {
    return s.conn.QueryContext(ctx, "SELECT short_key, url FROM storage WHERE user_id=$1", ctx.Value("user_id"));
}