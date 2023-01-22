package models

import "gopkg.in/guregu/null.v4"

type Category struct {
	Id   int         `json:"categoryId"`
	Name null.String `json:"name"`
	Timestamps
}
