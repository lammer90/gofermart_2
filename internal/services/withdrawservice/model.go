package withdrawservice

import (
	"errors"
	"github.com/lammer90/gofermart/internal/dto/withdraw"
)

type WithdrawService interface {
	Save(order, login string, sum float32) error
	FindAll(login string) ([]withdraw.WithdrawResponse, error)
	FindBalance(login string) (*withdraw.BalanceResponse, error)
}

var NotValidLuhnSum = errors.New("not valid luhn sum")
var NotEnoughMoney = errors.New("not enough money")
