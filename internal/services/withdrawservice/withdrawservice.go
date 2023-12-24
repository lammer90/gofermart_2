package withdrawservice

import (
	"database/sql"
	"github.com/EClaesson/go-luhn"
	"github.com/lammer90/gofermart/internal/dto/withdraw"
	"github.com/lammer90/gofermart/internal/repository/balance"
	"github.com/lammer90/gofermart/internal/repository/withdrawstorage"
	"time"
)

type withdrawServiceImpl struct {
	withdrawRepository withdrawstorage.WithdrawRepository
	balanceRepository  balance.BalanceRepository
	db                 *sql.DB
}

func New(withdrawRepository withdrawstorage.WithdrawRepository, balanceRepository balance.BalanceRepository, db *sql.DB) WithdrawService {
	return &withdrawServiceImpl{withdrawRepository: withdrawRepository, balanceRepository: balanceRepository, db: db}
}

func (w withdrawServiceImpl) Save(order, login string, sum float32) error {
	valid, err := luhn.IsValid(order)
	if err != nil || !valid {
		return NotValidLuhnSum
	}

	tx, err := w.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	balanceSum, err := w.balanceRepository.FindByUser(login, tx)
	if err != nil {
		return err
	}
	if balanceSum.BonusSum < sum {
		return NotEnoughMoney
	}
	wit := &withdraw.Withdraw{Login: login, Order: order, Sum: sum, ProcessedAt: time.Now()}
	err = w.withdrawRepository.Save(wit, tx)
	if err != nil {
		return err
	}
	err = w.balanceRepository.WithdrawBonus(login, sum, tx)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (w withdrawServiceImpl) FindAll(login string) ([]withdraw.WithdrawResponse, error) {
	response := make([]withdraw.WithdrawResponse, 0)

	withdraws, err := w.withdrawRepository.FindByUser(login)
	if err != nil {
		return nil, err
	}
	if withdraws == nil {
		return nil, nil
	}

	for _, with := range withdraws {
		resp := withdraw.WithdrawResponse{
			Order:       with.Order,
			Sum:         with.Sum,
			ProcessedAt: with.ProcessedAt.Format("2006-01-02T15:04:05-07:00")}
		response = append(response, resp)
	}
	return response, nil
}

func (w withdrawServiceImpl) FindBalance(login string) (*withdraw.BalanceResponse, error) {
	balanceSum, err := w.balanceRepository.FindByUser(login, nil)
	if err != nil {
		return nil, err
	}
	return &withdraw.BalanceResponse{Balance: balanceSum.BonusSum, Withdrawn: balanceSum.WithdrawSum}, nil
}
