package orderservice

import (
	"errors"
	"github.com/lammer90/gofermart/internal/dto/order"
)

type OrderService interface {
	Save(number, login string) error
	FindAll(login string) ([]order.OrderResponse, error)
	FindAllToProcess() ([]order.Order, error)
	UpdateAccrual(login, number string, status string, accrual float32) error
}

var ErrOrderNumberAlreadyExistThisUser = errors.New("order number already exist this user")

var ErrOrderNumberAlreadyExistAnotherUser = errors.New("order number already exist another user")

var ErrNotValidLuhnSum = errors.New("not valid luhn sum")
