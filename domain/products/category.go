package products

import (
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/models"
	"gopkg.in/guregu/null.v4"
)

type Category struct {
	models.Category
}

func NewCategory(name string) (Category, error) {
	cat := Category{
		Category: models.Category{
			Name: null.NewString(name, true),
		},
	}

	return cat, cat.Validate()
}

func (cat *Category) Validate() error {
	errl := apperror.NewErrorList()
	errl.SetType(apperror.AnyRequired)
	errl.SetPrefix("body ValidationError: ")

	if cat.Name.String == "" {
		errl.Add(ErrCategoryNameEmpty)
	}

	return errl.BuildError()
}
