package orders

import (
	"context"
	"database/sql"
	"encoding/json"
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
	q := `SELECT
		p.id,
		p.name,
		p.stock,
		p.price,
		p.image_url,
		p.category_id,
		d.id AS disc_id,
		d.min_qty AS disc_min_qty,
		d.type AS disc_type,
		d.result AS disc_result,
		(CASE WHEN d.expired_at IS NOT NULL THEN FROM_UNIXTIME(d.expired_at, '%Y-%m-%dT%H:%i:%s.%fZ') ELSE NULL END) AS disc_expired_at,
		(CASE WHEN d.expired_at IS NOT NULL THEN FROM_UNIXTIME(d.expired_at, '%d %b %Y') ELSE NULL END) AS disc_expired_at_format,
		(
			CASE WHEN d.result IS NOT NULL AND d.type = 'PERCENT' THEN 
				CONCAT('Discount ', d.result, '%', ' Rp. ', FORMAT((p.price * d.min_qty) - (((p.price * d.min_qty) / 100) * d.result), 0, 'de_DE'))
			ELSE
				CONCAT('Buy ', d.min_qty,' only Rp. ',FORMAT(d.result, 0, 'de_DE'))
			END
		) AS disc_percent_string_format
	FROM products p
	LEFT JOIN discounts d ON d.product_id = p.id
	WHERE (JSON_SEARCH(?, 'one', p.id) IS NOT NULL)`

	prodIdList := make([]int, 0)
	for _, prod := range prods {
		prodIdList = append(prodIdList, prod.Id)
	}

	prodIdListJson, err := json.Marshal(prodIdList)
	if err != nil {
		orderRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	rows, err := orderRepo.db.QueryContext(ctx, q, prodIdListJson)
	if err != nil {
		orderRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	order := new(Order)
	iter := 0
	for rows.Next() {
		prod := prods[iter]
		prod.Discount = new(Discount)
		if err := rows.Scan(
			&prod.Id,
			&prod.Name,
			&prod.Stock,
			&prod.Price,
			&prod.ImageUrl,
			&prod.CategoryId,
			&prod.Discount.Id,
			&prod.Discount.MinQty,
			&prod.Discount.Type,
			&prod.Discount.Result,
			&prod.Discount.ExpiredAt,
			&prod.Discount.ExpiredAtFormat,
			&prod.Discount.StringFormat,
		); err != nil {
			if err != sql.ErrNoRows {
				orderRepo.log.Error(err.Error())
				return nil, commonerr.ErrInternal
			}

			break
		}

		iter++
		prod.CalcSubtotal()
		order.Subtotal += prod.TotalFinalPrice
	}

	order.Products = prods

	return order, nil
}

func (orderRepo *OrdersRepository) AddOrder(ctx context.Context, cashierId int, order *Order) (*Order, error) {
	// addOrder(in_cashier_id INT, in_receipt_id VARCHAR(5), in_order JSON)
	receiptId := fmt.Sprintf(
		"S%s%s",
		utils.RandStringWithLetters("1234567890", 3),
		utils.RandStringWithLetters("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 1),
	)

	// DELETE
	j, _ := json.Marshal(order)
	fmt.Printf("cashierId: %d, receiptId: %s, order: %s\n", cashierId, receiptId, j)
	q := `CALL addOrder(?,?,?)`
	row := orderRepo.db.QueryRowContext(ctx, q, cashierId, receiptId, order)
	resOrder := new(Order)
	if err := row.Scan(resOrder); err != nil {
		orderRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	return resOrder, nil
}
