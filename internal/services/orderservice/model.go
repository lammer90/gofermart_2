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

var OrderNumberAlreadyExistThisUser = errors.New("order number already exist this user")

var OrderNumberAlreadyExistAnotherUser = errors.New("order number already exist another user")

var NotValidLuhnSum = errors.New("not valid luhn sum")
