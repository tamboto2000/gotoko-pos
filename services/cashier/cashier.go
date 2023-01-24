package cashier

import (
	"context"
	"log"

	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"go.uber.org/zap"
)

type CashierServiceConfig func(cashierSvc *CashierService) error

type CashierService struct {
	cashierRepo  *cashiers.CashiersRepository
	cashierMRepo *cashiers.CashiersMemRepository
	log          *zap.Logger
}

func NewCashierService(cfgs ...CashierServiceConfig) (*CashierService, error) {
	svc := new(CashierService)
	for _, cfg := range cfgs {
		if err := cfg(svc); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
	}

	return svc, nil
}

func WithCashierRepository(cashierRepo *cashiers.CashiersRepository) CashierServiceConfig {
	return func(cashierSvc *CashierService) error {
		cashierSvc.cashierRepo = cashierRepo

		return nil
	}
}

func WithCashierMemRepository(cashierRepo *cashiers.CashiersMemRepository) CashierServiceConfig {
	return func(cashierSvc *CashierService) error {
		cashierSvc.cashierMRepo = cashierRepo

		return nil
	}
}

func WithLogger(log *zap.Logger) CashierServiceConfig {
	return func(cashierSvc *CashierService) error {
		cashierSvc.log = log

		return nil
	}
}

func (cashierSvc *CashierService) CreateCashier(ctx context.Context, cashier *cashiers.Cashier) error {
	if err := cashier.Validate(); err != nil {
		return err
	}

	// encrypt passcode
	if err := cashier.EncryptPasscode(); err != nil {
		cashierSvc.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	if err := cashierSvc.cashierRepo.Create(ctx, cashier); err != nil {
		return err
	}

	return nil
}

func (cashierSvc *CashierService) GetDetail(ctx context.Context, id int) (*cashiers.Cashier, error) {
	return cashierSvc.cashierRepo.Get(ctx, id)
}

func (cashierSvc *CashierService) GetList(ctx context.Context, limit, skip int) (*cashiers.Cashiers, error) {
	return cashierSvc.cashierRepo.GetList(ctx, limit, skip)
}

func (cashierSvc *CashierService) Update(ctx context.Context, cashier *cashiers.Cashier) error {
	if err := cashier.ValidateForUpdate(); err != nil {
		return err
	}

	if cashier.Passcode != "" {
		if err := cashier.EncryptPasscode(); err != nil {
			cashierSvc.log.Error(err.Error())
			return commonerr.ErrInternal
		}
	}

	if err := cashierSvc.cashierRepo.Update(ctx, cashier); err != nil {
		return err
	}

	return nil
}

func (cashierSvc *CashierService) Delete(ctx context.Context, id int) error {
	return cashierSvc.cashierRepo.Delete(ctx, id)
}
