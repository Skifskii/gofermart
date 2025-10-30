package orders

import (
	"errors"
	"gophermart/internal/model"
	"gophermart/internal/repository"
	"sort"
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
	GetOrder(orderNum string) (model.Order, error)
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

	// Ищем заказ среди сохраненных
	ord, err := om.repo.GetOrder(orderNum)
	if err == nil {
		// Если ошибок нет, значит заказ уже загружен в систему
		if ord.UserLogin == userLogin {
			return ErrUserUploadedThisOrder
		}
		return ErrAnotherUserUploadedThisOrder
	} else {
		// Ошибку отсутствия заказа пропускаем, остальные возвращаем
		if !errors.Is(err, repository.ErrOrderNotFound) {
			return err
		}
	}

	// Получаем информацию о заказе
	order, err := om.infoGetter.GetOrderInfo(orderNum)
	if err != nil {
		return err
	}

	return om.repo.AddOrder(userLogin, order)
}

func validateOrderNumFormat(orderNum string) bool {
	return true // TODO:
}
