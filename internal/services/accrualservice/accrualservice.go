package accrualservice

import (
	"encoding/json"
	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/services/orderservice"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type accrualScheduledServiceImpl struct {
	orderService   orderservice.OrderService
	accrualAddress string
}

func New(orderService orderservice.OrderService, accrualAddress string) AccrualScheduledService {
	return &accrualScheduledServiceImpl{orderService: orderService, accrualAddress: accrualAddress}
}

func (a accrualScheduledServiceImpl) Start() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		logger.Log.Info(">> Start")
		numbers, err := a.orderService.FindAllToProcess()
		if err != nil {
			logger.Log.Error("Error during get orders to process", zap.Error(err))
			continue
		}
		for _, o := range numbers {
			response, err := http.Get(a.accrualAddress + "/api/orders/" + o.Number)
			if err != nil {
				logger.Log.Error("Error during get accrual by number "+o.Number, zap.Error(err))
				continue
			}
			var accrualResponse AccrualResponse
			dec := json.NewDecoder(response.Body)
			err = dec.Decode(&accrualResponse)
			if err != nil {
				logger.Log.Error("Error during get accrual by number "+o.Number, zap.Error(err))
				continue
			}
			a.orderService.UpdateAccrual(o.Login, o.Number, accrualResponse.Status, accrualResponse.Accrual)
			response.Body.Close()
		}
	}
}
