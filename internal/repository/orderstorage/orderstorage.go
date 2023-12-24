package orderstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lammer90/gofermart/internal/dto/order"
)

type dbOrderStorage struct {
	db *sql.DB
}

func New(db *sql.DB) OrderRepository {
	initDB(db)
	return &dbOrderStorage{db: db}
}

func (d dbOrderStorage) Save(order *order.Order) error {
	_, err := d.db.ExecContext(context.Background(), `
        INSERT INTO orders
        (order_number, login, status, accrual, uploaded_at)
        VALUES
        ($1, $2, $3, $4, $5);
    `, order.Number, order.Login, order.Status, order.Accrual, order.UploadedAt)
	return err
}

func (d dbOrderStorage) FindByUser(login string) ([]order.Order, error) {
	rows, err := d.db.QueryContext(context.Background(), `
        SELECT
            o.order_number,
            o.login,
            o.status,
            o.accrual,
            o.uploaded_at
        FROM orders o
        WHERE
            o.login = $1
        ORDER BY o.uploaded_at
    `, login)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]order.Order, 0)
	for rows.Next() {
		var o order.Order
		err = rows.Scan(&o.Number, &o.Login, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d dbOrderStorage) FindByNumber(number string) (*order.Order, error) {
	rows, err := d.db.QueryContext(context.Background(), `
        SELECT
            o.order_number,
            o.login,
            o.status,
            o.accrual,
            o.uploaded_at
        FROM orders o
        WHERE
            o.order_number = $1
    `, number)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var noRows = true
	var o order.Order
	for rows.Next() {
		noRows = false
		err = rows.Scan(&o.Number, &o.Login, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if noRows {
		return nil, nil
	}
	return &o, nil
}

func (d dbOrderStorage) FindNumbersToProcess() ([]order.Order, error) {
	rows, err := d.db.QueryContext(context.Background(), `
        SELECT o.order_number, o.login, o.status, o.accrual, o.uploaded_at FROM orders o WHERE o.status in (1,2)`)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]order.Order, 0)
	for rows.Next() {
		var o order.Order
		err = rows.Scan(&o.Number, &o.Login, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d dbOrderStorage) Update(order *order.Order, tx *sql.Tx) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(context.Background(), `
        UPDATE orders SET status = $1, accrual = $2 WHERE order_number = $3;
    `, order.Status, order.Accrual, order.Number)
	} else {
		_, err = d.db.ExecContext(context.Background(), `
        UPDATE orders SET status = $1, accrual = $2 WHERE order_number = $3;
    `, order.Status, order.Accrual, order.Number)
	}
	return err
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS orders (
            order_number varchar Primary Key,
            login varchar,
            status int,
            accrual numeric(15, 2),
            uploaded_at timestamp
        )
    `)
	db.ExecContext(ctx, `
        CREATE UNIQUE INDEX IF NOT EXISTS order_number_url_idx ON orders (order_number)
    `)
}
