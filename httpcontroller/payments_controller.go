package httpcontroller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/domain/payments"
	"github.com/tamboto2000/gotoko-pos/helpers/httpresponse"
	"github.com/tamboto2000/gotoko-pos/services/payment"
	"gopkg.in/guregu/null.v4"
)

type PaymentController struct {
	paySvc *payment.PaymentService
}

func NewPaymentController(paySvc *payment.PaymentService) *PaymentController {
	return &PaymentController{paySvc: paySvc}
}

func (payCtrl *PaymentController) CreatePayment(ctx *gin.Context) {
	payObj := new(payments.Payment)
	if err := ctx.BindJSON(payObj); err != nil {
		res := httpresponse.FromError(payObj.Validate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	if err := payCtrl.paySvc.CreatePayment(ctx, payObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(payObj)
	ctx.JSON(res.StatusCode, res)
}

func (payCtrl *PaymentController) GetPaymentDetail(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	paymentIdStr := ctx.Param("paymentId")
	paymentId, err := strconv.Atoi(paymentIdStr)
	if err != nil {
		errl.Add(apperror.New(`"paymentId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "paymentId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	pay, err := payCtrl.paySvc.GetPaymentDetail(ctx, paymentId)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(pay)
	ctx.JSON(res.StatusCode, res)
}

func (payCtrl *PaymentController) UpdatePayment(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	paymentIdStr := ctx.Param("paymentId")
	paymentId, err := strconv.Atoi(paymentIdStr)
	if err != nil {
		errl.Add(apperror.New(`"paymentId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "paymentId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	payObj := new(payments.Payment)
	if err := ctx.BindJSON(payObj); err != nil {
		res := httpresponse.FromError(payObj.ValidateForUpdate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	payObj.Id = null.NewInt(int64(paymentId), true)
	if err := payCtrl.paySvc.UpdatePayment(ctx, payObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}

func (payCtrl *PaymentController) DeletePayment(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	paymentIdStr := ctx.Param("paymentId")
	paymentId, err := strconv.Atoi(paymentIdStr)
	if err != nil {
		errl.Add(apperror.New(`"paymentId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "paymentId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	if err := payCtrl.paySvc.DeletePayment(ctx, paymentId); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}
