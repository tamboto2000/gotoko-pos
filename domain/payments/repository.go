package payments

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/guregu/null.v4"

	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
)

const (
	paymentTable = "payment_methods"
)

type PaymentsRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewPaymentRepository(db *sql.DB, log *zap.Logger) *PaymentsRepository {
	return &PaymentsRepository{db: db, log: log}
}

func (payRepo *PaymentsRepository) CreatePayment(ctx context.Context, pay *Payment) error {
	q := fmt.Sprintf(`INSERT INTO %s (name, type, logo_url) VALUES (?,?,?)`, paymentTable)
	res, err := payRepo.db.ExecContext(ctx, q, pay.Name, pay.Type, pay.LogoUrl)
	if err != nil {
		payRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		payRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q = fmt.Sprintf(`SELECT created_at, updated_at FROM %s WHERE id = ?`, paymentTable)
	row := payRepo.db.QueryRowContext(ctx, q, lastId)
	if err := row.Scan(
		&pay.CreatedAt,
		&pay.UpdatedAt,
	); err != nil {
		payRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	pay.Id = null.NewInt(lastId, true)

	return nil
}

func (payRepo *PaymentsRepository) GetPaymentDetail(ctx context.Context, id int) (*Payment, error) {
	q := fmt.Sprintf(`SELECT id, name, type, logo_url FROM %s WHERE id = ?`, paymentTable)
	row := payRepo.db.QueryRowContext(ctx, q, id)
	pay := new(Payment)
	if err := row.Scan(
		&pay.Id,
		&pay.Name,
		&pay.Type,
		&pay.LogoUrl,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPaymentNotFound
		}
	}

	return pay, nil
}

func (payRepo *PaymentsRepository) UpdatePayment(ctx context.Context, pay *Payment) error {
	q := fmt.Sprintf(`SELECT id FROM %s WHERE id = ?`, paymentTable)
	var id int
	row := payRepo.db.QueryRowContext(ctx, q, pay.Id)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrPaymentNotFound
		}

		payRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	/*
		updatePaymentMethod(
			in_id INT,
		    in_name VARCHAR(50),
		    in_type ENUM('CASH', 'E-WALLET', 'EDC'),
		    in_logo TEXT
		)
	*/
	q = `CALL updatePaymentMethod(?,?,?,?)`
	_, err := payRepo.db.ExecContext(ctx, q, id, pay.Name, pay.Type, pay.LogoUrl)
	if err != nil {
		payRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	return nil
}

func (payRepo *PaymentsRepository) DeletePayment(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, paymentTable)
	res, err := payRepo.db.ExecContext(ctx, q, id)
	if err != nil {
		payRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return ErrPaymentNotFound
	}

	return nil
}
