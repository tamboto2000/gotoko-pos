package products

type ProductList struct {
	Products []Product `json:"products"`
	Meta     Meta      `json:"meta"`
}

type Meta struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}
