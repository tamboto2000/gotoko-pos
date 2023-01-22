package products

type CategoryList struct {
	Categories []Category `json:"categories"`
	Meta       Meta       `json:"meta"`
}
