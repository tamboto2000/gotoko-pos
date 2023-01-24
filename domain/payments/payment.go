package payments

import (
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/models"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrNameEmpty       = apperror.New(`"name" is required`, apperror.AnyRequired, "name")
	ErrTypeEmpty       = apperror.New(`"type" is required`, apperror.AnyRequired, "type")
	ErrPaymentNotFound = apperror.New("Payment Not Found", apperror.NotFound, "")
)

type Payment struct {
	models.PaymentMethod
}

func NewPayment(name, t, logo string) (Payment, error) {
	pay := Payment{
		PaymentMethod: models.PaymentMethod{
			Name:    null.NewString(name, true),
			Type:    null.NewString(t, true),
			LogoUrl: null.NewString(logo, true),
		},
	}

	return pay, pay.Validate()
}

func (pay *Payment) Validate() error {
	errl := apperror.NewErrorList()
	errl.SetPrefix("body ValidationError: ")
	errl.SetType(apperror.AnyRequired)

	if pay.Name.String == "" {
		errl.Add(ErrNameEmpty)
	}

	if pay.Type.String == "" {
		errl.Add(ErrTypeEmpty)
	}

	return errl.BuildError()
}

func (pay *Payment) ValidateForUpdate() error {
	errl := apperror.NewErrorList()
	errl.SetPrefix("body ValidationError: ")
	errl.SetType(apperror.AnyRequired)

	if pay.Name.String == "" && pay.Type.String == "" &&
		pay.LogoUrl.String == "" {
		errl.Add(apperror.NewWithPeers(
			`"value" must contain at least one of [name, logo, type]`,
			apperror.ObjecMissing,
			[]string{},
			"value", []string{"name", "logo", "type"},
		))
	}

	return errl.BuildError()
}
