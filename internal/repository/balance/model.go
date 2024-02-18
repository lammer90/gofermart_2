package balance

import "database/sql"

type BalanceRepository interface {
	CreateBalance(login string) error
	AddBonus(login string, sumToAdd float32, tx *sql.Tx) error
	WithdrawBonus(login string, sumToMinus float32, tx *sql.Tx) error
	FindByUser(login string, tx *sql.Tx) (*Balance, error)
}

type Balance struct {
	Login       string
	BonusSum    float32
	WithdrawSum float32
}
