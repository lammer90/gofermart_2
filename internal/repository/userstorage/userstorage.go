package userstorage

import (
	"context"
	"database/sql"
	"errors"
)

type dbUserStorage struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	initDB(db)
	return &dbUserStorage{db: db}
}

func (d *dbUserStorage) Save(login, authHash string) error {
	_, err := d.db.ExecContext(context.Background(), `
        INSERT INTO users
        (login, auth_hash)
        VALUES
        ($1, $2);
    `, login, authHash)
	return err
}

func (d *dbUserStorage) Find(login string) (string, error) {
	row := d.db.QueryRowContext(context.Background(), `
        SELECT
            u.auth_hash
        FROM users u
        WHERE
            u.login = $1
    `, login)

	var authHash string
	err := row.Scan(&authHash)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	return authHash, nil
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS users (
            login varchar Primary Key,
            auth_hash varchar
        )
    `)
}
