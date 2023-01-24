package models

import (
	"gopkg.in/guregu/null.v4"
)

type PaymentMethod struct {
	Id      null.Int    `json:"paymentId"`
	Name    null.String `json:"name"`
	Type    null.String `json:"type"`
	LogoUrl null.String `json:"logoUrl"`
	Timestamps
}
