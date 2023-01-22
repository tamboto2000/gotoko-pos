package cashiers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/dbiface"
	"github.com/tamboto2000/gotoko-pos/models"
	"go.uber.org/zap"
)

const (
	cashiersTable    = "cashiers"
	cashierSessTable = "cashier_sessions"
)

type CashiersRepository struct {
	query dbiface.Query
	db    *sql.DB
	tx    *sql.Tx
	log   *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) *CashiersRepository {
	return &CashiersRepository{
		db:    db,
		query: db,
		log:   log,
	}
}

func (cashierRepo *CashiersRepository) NewTx(ctx context.Context) (*CashiersRepository, error) {
	tx, err := cashierRepo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &CashiersRepository{
		db:    cashierRepo.db,
		tx:    tx,
		query: tx,
		log:   cashierRepo.log,
	}, nil
}

func (cashierRepo *CashiersRepository) Get(ctx context.Context, id int) (*Cashier, error) {
	q := `SELECT id, name FROM ` + cashiersTable + ` WHERE id = ?`
	row := cashierRepo.query.QueryRowContext(ctx, q, id)
	cashier := new(Cashier)
	if err := row.Scan(
		&cashier.Id,
		&cashier.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCashierNotFound
		}

		cashierRepo.log.Error(err.Error())
		return nil, err
	}

	return cashier, nil
}

func (cashierRepo *CashiersRepository) GetList(ctx context.Context, limit, skip int) (*Cashiers, error) {
	q := fmt.Sprintf("SELECT id, name FROM %s ORDER BY id LIMIT ? OFFSET ?", cashiersTable)
	rows, err := cashierRepo.query.QueryContext(ctx, q, limit, skip)
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	cashiers := new(Cashiers)

	for rows.Next() {
		cashier := Cashier{}
		if err := rows.Scan(
			&cashier.Id,
			&cashier.Name,
		); err != nil {
			if err == sql.ErrNoRows {
				break
			}

			cashierRepo.log.Error(err.Error())
			return nil, commonerr.ErrInternal
		}

		cashiers.Cashiers = append(cashiers.Cashiers, cashier)
	}

	q = "SELECT COUNT(id) FROM " + cashiersTable
	row := cashierRepo.query.QueryRowContext(ctx, q)
	var count int
	if err := row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			cashierRepo.log.Error(err.Error())
			return nil, commonerr.ErrInternal
		}
	}

	cashiers.Meta.Total = count
	cashiers.Meta.Limit = limit
	cashiers.Meta.Skip = skip

	return cashiers, nil
}

func (cashierRepo *CashiersRepository) Create(ctx context.Context, c *Cashier) error {
	q := fmt.Sprintf(`INSERT INTO %s (name, passcode) VALUES (?,?)`, cashiersTable)
	res, err := cashierRepo.query.ExecContext(ctx, q, c.Name, c.encryptPasscode)
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q = fmt.Sprintf(`SELECT created_at, updated_at FROM %s WHERE id = LAST_INSERT_ID()`, cashiersTable)
	row := cashierRepo.query.QueryRowContext(ctx, q)
	c.CreatedAt = new(time.Time)
	c.UpdatedAt = new(time.Time)
	if err := row.Scan(
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	c.Id = int(lastId)

	return nil
}

func (cashierRepo *CashiersRepository) Update(ctx context.Context, c *Cashier) error {
	q := fmt.Sprintf(`SELECT id FROM %s WHERE id = ?`, cashiersTable)
	var id int
	row := cashierRepo.query.QueryRowContext(ctx, q, c.Id)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrCashierNotFound
		}

		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q = `SET @name = ?, @passcode = ?`
	_, err := cashierRepo.query.ExecContext(ctx, q, c.Name, c.encryptPasscode)
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q = `UPDATE cashiers SET
		name = CASE WHEN @name = '' THEN name ELSE @name END,
		passcode = CASE WHEN @passcode = '' THEN passcode ELSE @passcode END
	WHERE id = 4;`

	_, err = cashierRepo.query.ExecContext(ctx, q)
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	return nil
}

func (cashierRepo *CashiersRepository) GetPasscode(ctx context.Context, id int) (string, error) {
	q := fmt.Sprintf(`SELECT passcode FROM %s WHERE id = ?`, cashiersTable)
	row := cashierRepo.query.QueryRowContext(ctx, q, id)
	var pass string
	if err := row.Scan(&pass); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrCashierNotFound
		}

		cashierRepo.log.Error(err.Error())
		return "", commonerr.ErrInternal
	}

	return pass, nil
}

func (cashierRepo *CashiersRepository) CreateSession(ctx context.Context, c *models.CashierSession) error {
	q := fmt.Sprintf(`INSERT INTO %s (cashier_id, issued_at) VALUES (?,?)`, cashierSessTable)
	_, err := cashierRepo.query.ExecContext(ctx, q, c.CashierId, c.IssuedAt)
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	return nil
}

func (cashierRepo *CashiersRepository) DeleteSession(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE cashier_id = ?`, cashierSessTable)
	_, err := cashierRepo.query.ExecContext(ctx, q, id)
	if err != nil {
		cashierRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	return nil
}

func (cashierRepo *CashiersRepository) GetSession(ctx context.Context, cashierId int, issuedAt int64) (*models.CashierSession, error) {
	q := fmt.Sprintf(`SELECT cashier_id, issued_at FROM %s WHERE cashier_id = ? AND issued_at = ?`, cashierSessTable)
	row := cashierRepo.query.QueryRowContext(ctx, q, cashierId, issuedAt)

	sess := new(models.CashierSession)
	if err := row.Scan(
		&sess.CashierId,
		&sess.IssuedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCashierSessionNotFound
		}

		cashierRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	return sess, nil
}
