package payment

import (
	"context"

	"github.com/tamboto2000/gotoko-pos/domain/payments"
	"go.uber.org/zap"
)

type PaymentServiceConfig func(*PaymentService) error

type PaymentService struct {
	payRepo *payments.PaymentsRepository
	log     *zap.Logger
}

func NewPaymentService(cfgs ...PaymentServiceConfig) (*PaymentService, error) {
	paySvc := new(PaymentService)
	for _, cfg := range cfgs {
		if err := cfg(paySvc); err != nil {
			return nil, err
		}
	}

	return paySvc, nil
}

func WithPaymentsRepository(repo *payments.PaymentsRepository) PaymentServiceConfig {
	return func(ps *PaymentService) error {
		ps.payRepo = repo
		return nil
	}
}

func WithLogger(log *zap.Logger) PaymentServiceConfig {
	return func(ps *PaymentService) error {
		ps.log = log
		return nil
	}
}

func (paySvc *PaymentService) CreatePayment(ctx context.Context, pay *payments.Payment) error {
	if err := pay.Validate(); err != nil {
		return err
	}

	return paySvc.payRepo.CreatePayment(ctx, pay)
}

func (paySvc *PaymentService) GetPaymentDetail(ctx context.Context, id int) (*payments.Payment, error) {
	return paySvc.payRepo.GetPaymentDetail(ctx, id)
}

func (paySvc *PaymentService) UpdatePayment(ctx context.Context, pay *payments.Payment) error {
	if err := pay.ValidateForUpdate(); err != nil {
		return err
	}

	return paySvc.payRepo.UpdatePayment(ctx, pay)
}

func (paySvc *PaymentService) DeletePayment(ctx context.Context, id int) error {
	return paySvc.payRepo.DeletePayment(ctx, id)
}
