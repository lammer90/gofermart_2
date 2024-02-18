package withdrawstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lammer90/gofermart/internal/dto/withdraw"
)

type dbWithdrawStorage struct {
	db *sql.DB
}

func New(db *sql.DB) WithdrawRepository {
	initDB(db)
	return &dbWithdrawStorage{db: db}
}

func (d dbWithdrawStorage) Save(withdraw *withdraw.Withdraw, tx *sql.Tx) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(context.Background(), `
        INSERT INTO withdraws
        (withdraw_order, login, withdraw_sum, processed_at)
        VALUES
        ($1, $2, $3, $4);
    `, withdraw.Order, withdraw.Login, withdraw.Sum, withdraw.ProcessedAt)
	} else {
		_, err = d.db.ExecContext(context.Background(), `
        INSERT INTO withdraws
        (withdraw_order, login, withdraw_sum, processed_at)
        VALUES
        ($1, $2, $3, $4);
    `, withdraw.Order, withdraw.Login, withdraw.Sum, withdraw.ProcessedAt)
	}
	return err
}

func (d dbWithdrawStorage) FindByUser(login string) ([]withdraw.Withdraw, error) {
	rows, err := d.db.QueryContext(context.Background(), `
        SELECT
            w.withdraw_order,
            w.login,
            w.withdraw_sum,
            w.processed_at
        FROM withdraws w
        WHERE
            w.login = $1
        ORDER BY w.processed_at
    `, login)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]withdraw.Withdraw, 0)
	for rows.Next() {
		var w withdraw.Withdraw
		err = rows.Scan(&w.Order, &w.Login, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, w)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d dbWithdrawStorage) FindSumByUser(login string) (float32, error) {
	rows, err := d.db.QueryContext(context.Background(), `
        SELECT sum(w.withdraw_sum) FROM withdraws w WHERE w.login = $1
    `, login)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var s float32
	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			return 0, err
		}
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}
	return s, nil
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS withdraws (
            withdraw_order varchar,
            login varchar,
            withdraw_sum numeric(15, 2),
            processed_at timestamp
        )
    `)
	db.ExecContext(ctx, `
        CREATE INDEX IF NOT EXISTS login_withdraws_idx ON withdraws (login)
    `)
}
