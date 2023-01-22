package models

type Order struct {
	Id          int
	CashierId   int
	TotalPrice  uint
	TotalPaid   uint
	TotalReturn uint
	ReceiptId   string
	Timestamps
}
