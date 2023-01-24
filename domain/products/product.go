package products

import (
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/models"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrCategoryIdEmpty = apperror.New(`"categoryId" is required`, apperror.AnyRequired, "categoryId")
	ErrNameEmpty       = apperror.New(`"name" is required`, apperror.AnyRequired, "name")
	ErrImageEmpty      = apperror.New(`"image" is required`, apperror.AnyRequired, "image")
	ErrPriceEmpty      = apperror.New(`"price" is required`, apperror.AnyRequired, "price")
	ErrStockEmpty      = apperror.New(`"stock" is required`, apperror.AnyRequired, "stock")

	ErrDiscountQtyEmpty       = apperror.New(`"discount.qty" is required`, apperror.AnyRequired, "discount.qty")
	ErrDiscountTypeEmpty      = apperror.New(`"discount.type" is required`, apperror.AnyRequired, "discount.type")
	ErrDiscountResultEmpty    = apperror.New(`"discount.result" is required`, apperror.AnyRequired, "discount.result")
	ErrDiscountExpiredAtEmpty = apperror.New(`"discount.expiredAt" is required`, apperror.AnyRequired, "discount.expiredAt")
	ErrCategoryNameEmpty      = apperror.New(`"name" is required`, apperror.AnyRequired, "name")

	ErrProductNotFound  = apperror.New("Product Not Found", apperror.NotFound, "")
	ErrCategoryNotFound = apperror.New("Category Not Found", apperror.NotFound, "")
)

const (
	DiscountPercent = "PERCENT"
	DicountBuyN     = "BUY_N "
)

type Product struct {
	models.Product
	Category *Category        `json:"category,omitempty"`
	Discount *models.Discount `json:"discount,omitempty"`
}

func NewProduct(
	categoryId int,
	name string,
	imageUrl string,
	price uint,
	stock uint,
) (Product, error) {
	prod := Product{
		Product: models.Product{
			CategoryId: categoryId,
			Name:       name,
			ImageUrl:   imageUrl,
			Price:      null.NewInt(int64(price), true),
			Stock:      null.NewInt(int64(stock), true),
		},
	}

	return prod, prod.Validate()
}

func (prod *Product) Validate() error {
	errl := apperror.NewErrorList()
	errl.SetType(apperror.AnyRequired)
	errl.SetPrefix("body ValidationError: ")

	if prod.CategoryId <= 0 {
		errl.Add(ErrCategoryIdEmpty)
	}

	if prod.Name == "" {
		errl.Add(ErrNameEmpty)
	}

	if prod.ImageUrl == "" {
		errl.Add(ErrImageEmpty)
	}

	if !prod.Price.Valid {
		errl.Add(ErrPriceEmpty)
	}

	if !prod.Stock.Valid {
		errl.Add(ErrStockEmpty)
	}

	if prod.Discount != nil {
		if prod.Discount.MinQty.Int64 <= 0 {
			errl.Add(ErrDiscountQtyEmpty)
		}

		if prod.Discount.Type.String == "" {
			errl.Add(ErrDiscountTypeEmpty)
		}

		if prod.Discount.Result.Int64 <= 0 {
			errl.Add(ErrDiscountResultEmpty)
		}

		if prod.Discount.ExpiredAt.Int64 <= 0 {
			errl.Add(ErrDiscountExpiredAtEmpty)
		}
	}

	return errl.BuildError()
}

func (prod *Product) ValidateForUpdate() error {
	errl := apperror.NewErrorList()
	errl.SetPrefix("body ValidationError: ")
	errl.SetType(apperror.ObjecMissing)

	if prod.CategoryId == 0 && prod.Name == "" &&
		prod.ImageUrl == "" && !prod.Price.Valid &&
		!prod.Stock.Valid {
		errl.Add(apperror.NewWithPeers(
			`"value" must contain at least one of [categoryId, name, image, price, stock]`,
			apperror.ObjecMissing,
			[]string{},
			"value",
			[]string{
				"categoryId",
				"name",
				"image",
				"price",
				"stock",
			},
		))
	}

	return errl.BuildError()
}

func (prod *Product) SetDiscount(qty uint, ty string, result uint, expired int64) error {
	prod.Discount = &models.Discount{
		MinQty:    null.NewInt(int64(qty), true),
		Type:      null.NewString(ty, true),
		Result:    null.NewInt(int64(result), true),
		ExpiredAt: null.NewInt(expired, true),
	}

	return prod.Validate()
}
