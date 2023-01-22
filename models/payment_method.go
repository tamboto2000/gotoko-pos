package models

import (
	"gopkg.in/guregu/null.v4"
)

type PaymentMethod struct {
	Id      int
	Name    string
	Type    string
	LogoUrl null.String
	Timestamps
}
