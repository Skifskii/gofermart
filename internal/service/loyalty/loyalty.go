package loyalty

import (
	"errors"
	"slices"
)

var ErrWrongOrderNumFormat = errors.New("wrong order number format")
var ErrUserUploadedThisOrder = errors.New("the order number has already been uploaded by this user")
var ErrAnotherUserUploadedThisOrder = errors.New("the order number has already been uploaded by another user")

var ErrOrderNotRegistered = errors.New("the order is not registered in the payment system")
var ErrRequestLimitReached = errors.New("the number of requests to the service has been exceeded")

type LoyaltyPointsSystem struct {
	repo Repository
}

type Repository interface {
	GetOrderNums(login string) (orderNums []string, err error)
	AddOrderNum(login, orderNum string) error
}

type Status string

const (
	StatusRegistered Status = "REGISTERED"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

type Order struct {
	Num     string  `json:"order"`
	Status  Status  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func New(repo Repository) *LoyaltyPointsSystem {
	return &LoyaltyPointsSystem{repo: repo}
}

func (ls *LoyaltyPointsSystem) UploadOrderNum(userLogin, orderNum string) error {
	// Проверяем формат номера заказа
	if ok := validateOrderNumFormat(orderNum); !ok {
		return ErrWrongOrderNumFormat
	}

	// Проверяем, не загружал ли пользователь этот заказ ранее
	userOrderNums, err := ls.repo.GetOrderNums(userLogin)
	if err != nil {
		return err
	}
	if slices.Contains(userOrderNums, orderNum) {
		return ErrUserUploadedThisOrder
	}

	// Получаем информацию о заказе
	order, err := getInfoByOrderNum(orderNum)
	if err != nil {
		if errors.Is(err, ErrOrderNotRegistered) {
			return ErrWrongOrderNumFormat // TODO: точно ли тут эта ошибка? По заданию непонятно, эта подходит лучше всего
		}
		return err
	}

	// Если заказ новый, записываем его в заказы пользователя
	if order.Status == StatusRegistered {
		ls.repo.AddOrderNum(userLogin, orderNum)
		return nil
	}

	// Если заказ не в начальном статусе, значит его уже загружал ДРУГОЙ пользователь.
	// Сценарий загрузки ЭТИМ пользователем уже проверили ранее
	return ErrAnotherUserUploadedThisOrder
}

func validateOrderNumFormat(orderNum string) bool {
	return true // TODO:
}

func getInfoByOrderNum(orderNum string) (Order, error) {
	// TODO: прикрутить сервис
	mockOrders := []Order{
		{
			Num:     "123",
			Status:  "PROCESSED",
			Accrual: 500,
		},
		{
			Num:     "456",
			Status:  "REGISTERED",
			Accrual: 500,
		},
	}

	for _, order := range mockOrders {
		if order.Num == orderNum {
			return order, nil
		}
	}

	return Order{}, ErrOrderNotRegistered
}
