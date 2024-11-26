package pg

import (
    "context"
    "database/sql"
)

type Storage struct {
    conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
    return &Storage{conn: conn}
}

func (s Storage) Bootstrap(ctx context.Context) error {
    tx, err := s.conn.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer tx.Rollback()

    tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS storage(
			short_key varchar(128),
			url TEXT
		)
    `)
    tx.ExecContext(ctx, `CREATE UNIQUE INDEX short_key_idx ON users (short_key)`)

    return tx.Commit()
}

func (s Storage) Ping() error {
    return s.conn.Ping()
}

func (s Storage) Insert(ctx context.Context, shortKey string, url string) {
	s.conn.ExecContext(ctx, "INSERT INTO storage (short_key, url) VALUES ($1, $2)", shortKey, url)
}

func (s Storage) Get(ctx context.Context, shortKey string) (redirectLink string, err error) {
	row := s.conn.QueryRowContext(ctx, "SELECT url FROM storage WHERE short_key=$1", shortKey)
	err = row.Scan(&redirectLink)
	return
}