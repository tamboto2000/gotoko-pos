package httpcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/domain/orders"
	"github.com/tamboto2000/gotoko-pos/helpers/httpresponse"
	"github.com/tamboto2000/gotoko-pos/services/order"
)

type OrdersController struct {
	orderSvc *order.OrderService
}

func NewOrdersController(orderSvc *order.OrderService) *OrdersController {
	return &OrdersController{orderSvc: orderSvc}
}

func (orderCtrl *OrdersController) GetOrderSubtotal(ctx *gin.Context) {
	ol := new(orders.OrderItemList)
	if err := ctx.BindJSON(ol); err != nil {
		errl := apperror.NewErrorList()
		errl.SetPrefix("param ValidationError: ")
		errl.SetType(apperror.ArrayBase)

		errl.Add(apperror.NewMinimal(`"value" must be an array`, apperror.ArrayBase, "value"))
		res := httpresponse.FromError(errl.BuildError())
		ctx.JSON(res.StatusCode, res)

		return
	}

	orderObj, err := orderCtrl.orderSvc.GetOrderSubtotal(ctx, *ol)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(orderObj)
	ctx.JSON(res.StatusCode, res)
}

func (orderCtrl *OrdersController) AddOrder(ctx *gin.Context) {
	orderObj := new(orders.Order)
	if err := ctx.BindJSON(orderObj); err != nil {
		res := httpresponse.FromError(orderObj.Validate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	newOrderObj, err := orderCtrl.orderSvc.AddOrder(ctx, orderObj)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(newOrderObj)
	ctx.JSON(res.StatusCode, res)
}
