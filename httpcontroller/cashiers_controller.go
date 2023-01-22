package httpcontroller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"github.com/tamboto2000/gotoko-pos/helpers/httpresponse"
	"github.com/tamboto2000/gotoko-pos/services/cashier"
)

type CashiersController struct {
	cashierSvc *cashier.CashierService
}

func NewCashierController(cashierSvc *cashier.CashierService) *CashiersController {
	return &CashiersController{
		cashierSvc: cashierSvc,
	}
}

func (cashierCtrl *CashiersController) Create(ctx *gin.Context) {
	cashierObj := new(cashiers.Cashier)
	if err := ctx.BindJSON(cashierObj); err != nil {
		res := httpresponse.FromError(new(cashiers.Cashier).Validate())
		ctx.JSON(res.StatusCode, res)
		return
	}

	if err := cashierCtrl.cashierSvc.CreateCashier(ctx, cashierObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)
		return
	}

	res := httpresponse.Success(cashierObj)
	ctx.JSON(res.StatusCode, res)
}

func (cashierCtrl *CashiersController) GetDetail(ctx *gin.Context) {
	cashierIdStr := ctx.Param("cashierId")
	cashierId, err := strconv.Atoi(cashierIdStr)
	if err != nil {
		res := httpresponse.FromError(apperror.New(
			`"cashierId" invalid, valid input are numeric (0-9)`,
			apperror.AnyRequired,
			"cashierId",
		))

		ctx.JSON(res.StatusCode, res)
		return
	}

	cashierObj, err := cashierCtrl.cashierSvc.GetDetail(ctx, cashierId)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)
		return
	}

	res := httpresponse.Success(cashierObj)
	ctx.JSON(res.StatusCode, res)
}

func (cashierCtrl *CashiersController) GetList(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	skipStr := ctx.Query("skip")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		res := httpresponse.FromError(apperror.New(
			`"limit" is invalid, valid input are numeric (0-9)`,
			apperror.AnyRequired,
			"limit",
		))

		ctx.JSON(res.StatusCode, res)
		return
	}

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		res := httpresponse.FromError(apperror.New(`"skip" is invalid, valid input are numeric (0-9)`, apperror.AnyRequired, "skip"))
		ctx.JSON(res.StatusCode, res)

		return
	}

	cashierList, err := cashierCtrl.cashierSvc.GetList(ctx, limit, skip)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(cashierList)
	ctx.JSON(res.StatusCode, res)
}

func (cashierCtrl *CashiersController) Update(ctx *gin.Context) {
	cashierIdStr := ctx.Param("cashierId")
	cashierId, err := strconv.Atoi(cashierIdStr)
	if err != nil {
		res := httpresponse.FromError(apperror.New(
			`"cashierId" invalid, valid input are numeric (0-9)`,
			apperror.AnyRequired,
			"cashierId",
		))

		ctx.JSON(res.StatusCode, res)
		return
	}

	cashierObj := new(cashiers.Cashier)
	if err := ctx.BindJSON(cashierObj); err != nil {
		res := httpresponse.FromError(new(cashiers.Cashier).ValidateForUpdate())
		ctx.JSON(res.StatusCode, res)
		return
	}

	cashierObj.Id = cashierId
	if err := cashierCtrl.cashierSvc.Update(ctx, cashierObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}
