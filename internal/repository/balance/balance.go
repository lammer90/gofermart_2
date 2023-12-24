package balance

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type dbBalanceStorage struct {
	db *sql.DB
}

func New(db *sql.DB) BalanceRepository {
	initDB(db)
	return &dbBalanceStorage{db: db}
}

func (d dbBalanceStorage) CreateBalance(login string) error {
	_, insErr := d.db.ExecContext(context.Background(), `INSERT INTO balance (login, balance_sum, withdraw_sum) VALUES($1, $2, $3);`, login, 0, 0)
	if insErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(insErr, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil
		}
	}
	return insErr
}

func (d dbBalanceStorage) AddBonus(login string, sumToAdd float32, tx *sql.Tx) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(context.Background(), `
        UPDATE balance SET balance_sum = balance_sum + $1 WHERE login = $2
    `, sumToAdd, login)
	} else {
		_, err = d.db.ExecContext(context.Background(), `
        UPDATE balance SET balance_sum = balance_sum + $1 WHERE login = $2
    `, sumToAdd, login)
	}
	return err
}

func (d dbBalanceStorage) WithdrawBonus(login string, sumToMinus float32, tx *sql.Tx) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(context.Background(), `
        UPDATE balance SET balance_sum = balance_sum - $1, withdraw_sum = withdraw_sum + $2 WHERE login = $3
    `, sumToMinus, sumToMinus, login)
	} else {
		_, err = d.db.ExecContext(context.Background(), `
        UPDATE balance SET balance_sum = balance_sum - $1, withdraw_sum = withdraw_sum + $2 WHERE login = $3
    `, sumToMinus, sumToMinus, login)
	}
	return err
}

func (d dbBalanceStorage) FindByUser(login string, tx *sql.Tx) (*Balance, error) {
	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryContext(context.Background(), `
        SELECT
            b.login,
            b.balance_sum,
            b.withdraw_sum
        FROM balance b
        WHERE
            b.login = $1
        FOR UPDATE
    `, login)
	} else {
		rows, err = d.db.QueryContext(context.Background(), `
        SELECT
            b.login,
            b.balance_sum,
            b.withdraw_sum
        FROM balance b
        WHERE
            b.login = $1
    `, login)
	}

	if err != nil {
		return nil, err
	}

	var result Balance
	for rows.Next() {
		err = rows.Scan(&result.Login, &result.BonusSum, &result.WithdrawSum)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS balance (
            login varchar,
            balance_sum numeric(15, 2),
            withdraw_sum numeric(15, 2)
        )
    `)
	db.ExecContext(ctx, `
        CREATE UNIQUE INDEX IF NOT EXISTS login_balance_idx ON balance (login)
    `)
}
