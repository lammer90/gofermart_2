package withdraw

import "time"

type Withdraw struct {
	Login       string
	Order       string
	Sum         float32
	ProcessedAt time.Time
}

type WithdrawResponse struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type BalanceResponse struct {
	Balance   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
