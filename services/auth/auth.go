package auth

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tamboto2000/gotoko-pos/apperror"
	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"github.com/tamboto2000/gotoko-pos/helpers/aescrypt"
	"github.com/tamboto2000/gotoko-pos/helpers/jwtparse"
	"github.com/tamboto2000/gotoko-pos/models"
	"go.uber.org/zap"
)

var (
	ErrInvalidPass = apperror.New("Passcode Not Match", apperror.InvalidAuth, "")
)

type AuthServiceConfig func(*AuthService) error

type AuthService struct {
	cashierRepo  *cashiers.CashiersRepository
	cashierMRepo *cashiers.CashiersMemRepository
	log          *zap.Logger
}

func NewAuthService(cfgs ...AuthServiceConfig) (*AuthService, error) {
	authSvc := new(AuthService)

	for _, cfg := range cfgs {
		if err := cfg(authSvc); err != nil {
			return nil, err
		}
	}

	return authSvc, nil
}

func WithCashierRepository(repo *cashiers.CashiersRepository) AuthServiceConfig {
	return func(authSvc *AuthService) error {
		authSvc.cashierRepo = repo
		return nil
	}
}

func WithCashierMemRepository(repo *cashiers.CashiersMemRepository) AuthServiceConfig {
	return func(authSvc *AuthService) error {
		authSvc.cashierMRepo = repo
		return nil
	}
}

func WithLogger(log *zap.Logger) AuthServiceConfig {
	return func(authSvc *AuthService) error {
		authSvc.log = log
		return nil
	}
}

func (authSvc *AuthService) GetCashierPasscode(ctx context.Context, id int) (*cashiers.Cashier, error) {
	pass, err := authSvc.cashierRepo.GetPasscode(ctx, id)
	if err != nil {
		return nil, err
	}

	raw, err := hex.DecodeString(pass)
	if err != nil {
		authSvc.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	decPass, err := aescrypt.Decrypt([]byte(raw))
	if err != nil {
		authSvc.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	cashierObj := new(cashiers.Cashier)
	cashierObj.Passcode = string(decPass)
	return cashierObj, nil
}

func (authSvc *AuthService) CashierLogin(ctx context.Context, cashier *cashiers.Cashier) (map[string]interface{}, error) {
	if err := cashier.ValidateForLogin(); err != nil {
		return nil, err
	}

	pass, err := authSvc.cashierRepo.GetPasscode(ctx, cashier.Id)
	if err != nil {
		return nil, err
	}

	raw, err := hex.DecodeString(pass)
	if err != nil {
		authSvc.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	plainPass, err := aescrypt.Decrypt(raw)
	if err != nil {
		authSvc.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	if string(plainPass) != cashier.Passcode {
		return nil, ErrInvalidPass
	}

	iat := jwt.NewNumericDate(time.Now())
	token, err := jwtparse.BuildJWT(jwt.RegisteredClaims{
		IssuedAt: iat,
		Subject:  strconv.Itoa(cashier.Id),
	})

	if err != nil {
		authSvc.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	if err := authSvc.cashierRepo.CreateSession(ctx, &models.CashierSession{
		CashierId: cashier.Id,
		IssuedAt:  iat.Unix(),
	}); err != nil {
		return nil, err
	}

	authSvc.cashierMRepo.CreateSession(ctx, &models.CashierSession{
		CashierId: cashier.Id,
		IssuedAt:  iat.Unix(),
	})

	return map[string]interface{}{"token": token}, nil
}

func (authSvc *AuthService) CashierLogout(ctx context.Context, cashier *cashiers.Cashier) error {
	if err := cashier.ValidateForLogin(); err != nil {
		return err
	}

	pass, err := authSvc.cashierRepo.GetPasscode(ctx, cashier.Id)
	if err != nil {
		return err
	}

	raw, err := hex.DecodeString(pass)
	if err != nil {
		authSvc.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	plainPass, err := aescrypt.Decrypt(raw)
	if err != nil {
		authSvc.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	if string(plainPass) != cashier.Passcode {
		return ErrInvalidPass
	}

	if err := authSvc.cashierRepo.DeleteSession(ctx, cashier.Id); err != nil {
		return err
	}

	authSvc.cashierMRepo.DeleteSession(ctx, cashier.Id)

	return nil
}
