package orders

import (
	"errors"
	"gophermart/internal/model"
	"sort"
	"time"
)

var ErrWrongOrderNumFormat = errors.New("wrong order number format")
var ErrUserUploadedThisOrder = errors.New("the order number has already been uploaded by this user")
var ErrAnotherUserUploadedThisOrder = errors.New("the order number has already been uploaded by another user")

type OrdersManager struct {
	repo       Repository
	infoGetter OrderInfoGetter
}

type Repository interface {
	GetOrders(userLogin string) ([]model.Order, error)
	AddOrder(login string, order model.Order) error
}

type OrderInfoGetter interface {
	GetOrderInfo(orderNum string) (model.Order, error)
}

func New(repo Repository, infoGetter OrderInfoGetter) *OrdersManager {
	return &OrdersManager{
		repo:       repo,
		infoGetter: infoGetter,
	}
}

func (om *OrdersManager) GetOrders(userLogin string) ([]model.Order, error) {
	orders, err := om.repo.GetOrders(userLogin)
	if err != nil {
		return nil, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.After(orders[j].UploadedAt)
	})

	return orders, nil
}

func (om *OrdersManager) AddOrder(userLogin, orderNum string) error {
	// Проверяем формат номера заказа
	if ok := validateOrderNumFormat(orderNum); !ok {
		return ErrWrongOrderNumFormat
	}

	// Проверяем, не загружал ли пользователь этот заказ ранее
	userOrders, err := om.GetOrders(userLogin)
	if err != nil {
		return err
	}
	for _, order := range userOrders {
		if order.Number == orderNum {
			return ErrUserUploadedThisOrder
		}
	}

	// Получаем информацию о заказе
	order, err := om.infoGetter.GetOrderInfo(orderNum)
	if err != nil {
		return err
	}

	// Если заказ новый, записываем его в заказы пользователя
	if order.Status == model.StatusRegistered {
		order.UploadedAt = time.Now()
		return om.repo.AddOrder(userLogin, order)
	}

	// Если заказ не в начальном статусе, значит его уже загружал ДРУГОЙ пользователь.
	// Сценарий загрузки ЭТИМ пользователем уже проверили ранее
	return ErrAnotherUserUploadedThisOrder
}

func validateOrderNumFormat(orderNum string) bool {
	return true // TODO:
}
