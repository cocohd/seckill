package services

import (
	"seckill/datamodels"
	"seckill/repositories"
)

type IOrderService interface {
	InsertOrder(order *datamodels.Order) (int64, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(order *datamodels.Order) error
	GetOrderByID(int64) (*datamodels.Order, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message) (int64, error)
}

type OrderService struct {
	orderRepositories repositories.IOrderRepository
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{repository}
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	return o.orderRepositories.Insert(order)
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	return o.orderRepositories.Delete(orderID)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderRepositories.Update(order)
}

func (o *OrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	return o.orderRepositories.SelectByKey(orderID)
}

func (o *OrderService) GetAllOrder() (orders []*datamodels.Order, err error) {
	return o.orderRepositories.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (orderMap map[int]map[string]string, err error) {
	return o.orderRepositories.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserID:      message.UserId,
		ProductID:   message.ProductId,
		OrderStatus: datamodels.OrderSuccess,
	}

	return o.orderRepositories.Insert(order)
}
