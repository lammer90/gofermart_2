package withdrawstorage

import (
	"database/sql"
	"github.com/lammer90/gofermart/internal/dto/withdraw"
)

type WithdrawRepository interface {
	Save(withdraw *withdraw.Withdraw, tx *sql.Tx) error
	FindByUser(login string) ([]withdraw.Withdraw, error)
	FindSumByUser(login string) (float32, error)
}
