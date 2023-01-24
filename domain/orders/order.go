package orders

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/models"
)

var (
	ErrPaymentIdEmpty = apperror.New(`"paymentId" is required`, apperror.AnyRequired, "paymentId")
	ErrTotalPaidEmpty = apperror.New(`"totalPaid" is required`, apperror.AnyRequired, "totalPaid")
	ErrProductsEmpty  = apperror.New(`"products" is required`, apperror.AnyRequired, "products")
)

type OrderMeta struct {
	OrderId        int       `json:"orderId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CashiersId     int       `json:"cashiersId"`
	TotalPaid      int       `json:"totalPaid"`
	TotalPrice     int       `json:"totalPrice"`
	TotalReturn    int       `json:"totalReturn"`
	ReceiptId      string    `json:"receiptId"`
	PaymentTypesId int       `json:"paymentTypesId"`
}

type Order struct {
	OrderMeta *OrderMeta    `json:"order,omitempty"`
	Subtotal  int           `json:"subtotal,omitempty"`
	Products  OrderItemList `json:"products"`
	TotalPaid int           `json:"totalPaid,omitempty"`
	PaymentId int           `json:"paymentId,omitempty"`
}

func (o *Order) Validate() error {
	errl := apperror.NewErrorList()
	errl.SetPrefix("body ValidationError: ")
	errl.SetType(apperror.AnyRequired)

	if o.PaymentId <= 0 {
		errl.Add(ErrPaymentIdEmpty)
	}

	if o.TotalPaid <= 0 {
		errl.Add(ErrTotalPaidEmpty)
	}

	if o.Products == nil || len(o.Products) == 0 {
		errl.Add(ErrProductsEmpty)
	}

	return errl.BuildError()
}

func (o *Order) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		b, ok = value.([]uint8)
		if !ok {
			return errors.New("type assertion to []byte and []uint8 failed")
		}
	}

	return json.Unmarshal(b, o)
}

func (o *Order) Value() (driver.Value, error) {
	return json.Marshal(o)
}

type Discount struct {
	Id              int    `json:"-"`
	ProductId       int    `json:"-"`
	MinQty          int    `json:"qty"`
	Type            string `json:"type"`
	Result          int    `json:"result"`
	ExpiredAt       string `json:"expiredAt"`
	StringFormat    string `json:"stringFormat"`
	ExpiredAtFormat string `json:"expiredAtFormat"`
}

type OrderItem struct {
	models.Product
	Discount         *Discount `json:"discount"`
	Qty              int       `json:"qty"`
	TotalNormalPrice int       `json:"totalNormalPrice"`
	TotalFinalPrice  int       `json:"totalFinalPrice"`
}

type OrderItemList []OrderItem

func (ol OrderItemList) Validate() error {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.ArrayBase)

	if len(ol) == 0 {
		errl.Add(apperror.NewMinimal(`"value" must be an array`, apperror.ArrayBase, "value"))
	}

	return errl.BuildError()
}

func (ol OrderItemList) Value() (driver.Value, error) {
	return json.Marshal(ol)
}
