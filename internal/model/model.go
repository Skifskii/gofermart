package model

import "time"

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Status string

const (
	StatusRegistered Status = "REGISTERED"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

type Order struct {
	Number     string
	Status     Status
	Accrual    *float64
	UploadedAt time.Time
}
