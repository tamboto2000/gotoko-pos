package models

type Cashier struct {
	Id       int    `json:"cashierId,omitempty"`
	Name     string `json:"name,omitempty"`
	Passcode string `json:"passcode,omitempty"`
	Timestamps
}
