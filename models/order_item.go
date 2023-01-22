package models

type OrderItem struct {
	Id        int
	OrderId   int
	ProductId int
	Qty       int
	Timestamps
}
