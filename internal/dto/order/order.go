package order

import "time"

type Order struct {
	Login      string
	Number     string
	Status     Status
	Accrual    float32
	UploadedAt time.Time
}

type OrderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type Status int

const (
	NEW Status = iota + 1
	PROCESSING
	INVALID
	PROCESSED
)

var Statuses = []string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}
