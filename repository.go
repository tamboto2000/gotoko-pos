package main

import (
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"github.com/tamboto2000/gotoko-pos/domain/orders"
	"github.com/tamboto2000/gotoko-pos/domain/payments"
	"github.com/tamboto2000/gotoko-pos/domain/products"
	"go.uber.org/zap"
)

var (
	cashiersRepo  *cashiers.CashiersRepository
	cashiersMRepo *cashiers.CashiersMemRepository
	prodRepo      *products.ProductsRepository
	payRepo       *payments.PaymentsRepository
	orderRepo     *orders.OrdersRepository
)

func buildRepositories(logging *zap.Logger) {
	db, err := buildDatabase()
	if err != nil {
		logging.Fatal(err.Error())
	}

	cashiersRepo = cashiers.New(db, logging)
	cashiersMRepo = cashiers.NewMemory()
	prodRepo = products.NewProductsRepository(db, logging)
	payRepo = payments.NewPaymentRepository(db, logging)
	orderRepo = orders.NewOrdersRepository(db, logging)
}
