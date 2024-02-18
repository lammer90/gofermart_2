package accrualservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/services/orderservice"
)

type accrualScheduledServiceImpl struct {
	orderService   orderservice.OrderService
	accrualAddress string
}

func New(orderService orderservice.OrderService, accrualAddress string) AccrualScheduledService {
	return &accrualScheduledServiceImpl{orderService: orderService, accrualAddress: accrualAddress}
}

func (a accrualScheduledServiceImpl) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			logger.Log.Info(">> Start")
			numbers, err := a.orderService.FindAllToProcess()
			if err != nil {
				logger.Log.Error("Error during get orders to process", err)
				continue
			}
			for _, o := range numbers {
				response, err := http.Get(a.accrualAddress + "/api/orders/" + o.Number)
				if err != nil {
					logger.Log.Error("Error during get accrual by number "+o.Number, err)
					continue
				}
				var accrualResponse AccrualResponse
				dec := json.NewDecoder(response.Body)
				err = dec.Decode(&accrualResponse)
				if err != nil {
					logger.Log.Error("Error during get accrual by number "+o.Number, err)
					continue
				}
				a.orderService.UpdateAccrual(o.Login, o.Number, accrualResponse.Status, accrualResponse.Accrual)
				response.Body.Close()
			}
		case <-ctx.Done():
			return
		}
	}
}
