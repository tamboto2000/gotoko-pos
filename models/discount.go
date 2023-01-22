package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v4"
)

type Discount struct {
	Id        null.Int    `json:"-"`
	ProductId null.Int    `json:"-"`
	MinQty    null.Int    `json:"qty"`
	Type      null.String `json:"type"`
	Result    null.Int    `json:"result"`
	ExpiredAt null.Int    `json:"expiredAt"`
	Timestamps
}

func (d *Discount) MarshalJSON() ([]byte, error) {
	if !d.Type.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(*d)
}
