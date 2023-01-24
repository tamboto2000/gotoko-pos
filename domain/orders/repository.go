package orders

import (
	"context"
	"database/sql"
	"fmt"

	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/utils"
	"go.uber.org/zap"
)

type OrdersRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewOrdersRepository(db *sql.DB, log *zap.Logger) *OrdersRepository {
	return &OrdersRepository{
		db:  db,
		log: log,
	}
}

func (orderRepo *OrdersRepository) GetSubTotal(ctx context.Context, prods OrderItemList) (*Order, error) {
	// orderSubTotal(in_items JSON)
	q := `SELECT orderSubTotal(?)`
	row := orderRepo.db.QueryRowContext(ctx, q, prods)
	order := new(Order)
	if err := row.Scan(order); err != nil {
		orderRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	return order, nil
}

func (orderRepo *OrdersRepository) AddOrder(ctx context.Context, cashierId int, order *Order) (*Order, error) {
	// addOrder(in_cashier_id INT, in_receipt_id VARCHAR(5), in_order JSON)
	receiptId := fmt.Sprintf(
		"S%s%s",
		utils.RandStringWithLetters("1234567890", 3),
		utils.RandStringWithLetters("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 1),
	)
	q := `CALL addOrder(?,?,?)`
	row := orderRepo.db.QueryRowContext(ctx, q, cashierId, receiptId, order)
	resOrder := new(Order)
	if err := row.Scan(resOrder); err != nil {
		orderRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	return resOrder, nil
}
