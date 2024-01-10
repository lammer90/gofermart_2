package accrualservice

import "context"

type AccrualScheduledService interface {
	Start(ctx context.Context)
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}
