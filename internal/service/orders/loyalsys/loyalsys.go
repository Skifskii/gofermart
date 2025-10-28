package loyalsys

import (
	"errors"
	"gophermart/internal/model"
	"time"
)

var ErrOrderNotRegistered = errors.New("the order is not registered in the payment system")
var ErrRequestLimitReached = errors.New("the number of requests to the service has been exceeded") // TODO: может возвращаться при запросе в смежный сервис

type LoyaltySystem struct{}

func New() *LoyaltySystem {
	return &LoyaltySystem{}
}

type orderResponse struct {
	Order   string       `json:"order"`
	Status  model.Status `json:"status"`
	Accrual *float64     `json:"accrual"`
}

func (or *orderResponse) toDomain() model.Order {
	return model.Order{
		Number:     or.Order,
		Status:     or.Status,
		Accrual:    or.Accrual,
		UploadedAt: time.Time{},
	}
}

func (ls *LoyaltySystem) GetOrderInfo(orderNum string) (model.Order, error) {
	f500 := 500.
	// TODO: прикрутить сервис
	mockOrders := []orderResponse{
		{
			Order:   "123",
			Status:  "PROCESSED",
			Accrual: &f500,
		},
		{
			Order:   "456",
			Status:  "REGISTERED",
			Accrual: &f500,
		},
		{
			Order:  "789",
			Status: "REGISTERED",
		},
	}

	for _, order := range mockOrders {
		if order.Order == orderNum {
			return order.toDomain(), nil
		}
	}

	return model.Order{}, ErrOrderNotRegistered
}
