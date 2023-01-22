package models

import "gopkg.in/guregu/null.v4"

type Product struct {
	Id         int      `json:"productId"`
	Sku        string   `json:"sku"`
	Name       string   `json:"name"`
	Stock      null.Int `json:"stock"`
	Price      null.Int `json:"price"`
	ImageUrl   string   `json:"image"`
	CategoryId int      `json:"categoryId,omitempty"`
	Timestamps
}
