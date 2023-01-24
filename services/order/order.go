package order

import (
	"context"

	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/domain/orders"
	"go.uber.org/zap"
)

type OrderServiceConfig func(*OrderService) error

type OrderService struct {
	orderRepo *orders.OrdersRepository
	log       *zap.Logger
}

func NewOrderService(cfgs ...OrderServiceConfig) (*OrderService, error) {
	orderSvc := new(OrderService)
	for _, cfg := range cfgs {
		if err := cfg(orderSvc); err != nil {
			return nil, err
		}
	}

	return orderSvc, nil
}

func WithOrdersRepository(repo *orders.OrdersRepository) OrderServiceConfig {
	return func(os *OrderService) error {
		os.orderRepo = repo

		return nil
	}
}

func WithLogger(log *zap.Logger) OrderServiceConfig {
	return func(os *OrderService) error {
		os.log = log

		return nil
	}
}

func (orderSvc *OrderService) GetOrderSubtotal(ctx context.Context, ol orders.OrderItemList) (*orders.Order, error) {
	if err := ol.Validate(); err != nil {
		return nil, err
	}

	return orderSvc.orderRepo.GetSubTotal(ctx, ol)
}

func (orderSvc *OrderService) AddOrder(ctx context.Context, or *orders.Order) (*orders.Order, error) {
	if err := or.Validate(); err != nil {
		return nil, err
	}

	cashierId, ok := ctx.Value("cashier-id").(int)
	if !ok {
		orderSvc.log.Error("can not cast cashier-id to int")
		return nil, commonerr.ErrInternal
	}

	return orderSvc.orderRepo.AddOrder(ctx, cashierId, or)
}
