package main

import (
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"github.com/tamboto2000/gotoko-pos/domain/products"
	"go.uber.org/zap"
)

var (
	cashiersRepo  *cashiers.CashiersRepository
	cashiersMRepo *cashiers.CashiersMemRepository
	prodRepo      *products.ProductsRepository
)

func buildRepositories(logging *zap.Logger) {
	db, err := buildDatabase()
	if err != nil {
		logging.Fatal(err.Error())
	}

	cashiersRepo = cashiers.New(db, logging)
	cashiersMRepo = cashiers.NewMemory()
	prodRepo = products.NewProductsRepository(db, logging)
}
