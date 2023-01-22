package models

type CashierSession struct {
	Id        int   `json:"id"`
	CashierId int   `json:"cashier_id"`
	IssuedAt  int64 `json:"issued_at"`
	Timestamps
}
