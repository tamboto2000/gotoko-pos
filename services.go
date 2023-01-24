package main

import (
	"github.com/tamboto2000/gotoko-pos/services/auth"
	"github.com/tamboto2000/gotoko-pos/services/cashier"
	"github.com/tamboto2000/gotoko-pos/services/order"
	"github.com/tamboto2000/gotoko-pos/services/payment"
	"github.com/tamboto2000/gotoko-pos/services/product"
	"go.uber.org/zap"
)

var (
	cashierSvc *cashier.CashierService
	authSvc    *auth.AuthService
	prodSvc    *product.ProductService
	paySvc     *payment.PaymentService
	orderSvc   *order.OrderService
)

// buildServices build services with their configuration.
// Please use this function to build your services
func buildServices(logging *zap.Logger) {
	var err error
	// CashierService
	cashierSvc, err = cashier.NewCashierService(
		cashier.WithCashierRepository(cashiersRepo),
		cashier.WithCashierMemRepository(cashiersMRepo),
		cashier.WithLogger(logging),
	)

	if err != nil {
		logging.Fatal(err.Error())
	}

	// AuthService
	authSvc, err = auth.NewAuthService(
		auth.WithCashierRepository(cashiersRepo),
		auth.WithCashierMemRepository(cashiersMRepo),
		auth.WithLogger(logging),
	)

	if err != nil {
		logging.Fatal(err.Error())
	}

	// ProductService
	prodSvc, err = product.NewProductService(
		product.WithProductsRepository(prodRepo),
		product.WithLogger(logging),
	)
	if err != nil {
		logging.Fatal(err.Error())
	}

	// PaymentService
	paySvc, err = payment.NewPaymentService(
		payment.WithPaymentsRepository(payRepo),
		payment.WithLogger(logging),
	)

	if err != nil {
		logging.Fatal(err.Error())
	}

	// OrderService
	orderSvc, err = order.NewOrderService(
		order.WithOrdersRepository(orderRepo),
		order.WithLogger(logging),
	)
}
