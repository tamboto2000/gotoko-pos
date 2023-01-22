package httpcontroller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"github.com/tamboto2000/gotoko-pos/helpers/httpresponse"
	"github.com/tamboto2000/gotoko-pos/services/auth"
)

type AuthController struct {
	authSvc *auth.AuthService
}

func NewAuthController(authSvc *auth.AuthService) *AuthController {
	return &AuthController{
		authSvc: authSvc,
	}
}

func (authCtrl *AuthController) GetCashierPasscode(ctx *gin.Context) {
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

	cashierObj, err := authCtrl.authSvc.GetCashierPasscode(ctx, cashierId)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(cashierObj)
	ctx.JSON(res.StatusCode, res)
}

func (authCtrl *AuthController) CashierLogin(ctx *gin.Context) {
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
		res := httpresponse.FromError(new(cashiers.Cashier).ValidateForLogin())
		ctx.JSON(res.StatusCode, res)
		return
	}

	cashierObj.Id = cashierId
	token, err := authCtrl.authSvc.CashierLogin(ctx, cashierObj)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(token)
	ctx.JSON(res.StatusCode, res)
}

func (authCtrl *AuthController) CashierLogout(ctx *gin.Context) {
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
		res := httpresponse.FromError(new(cashiers.Cashier).ValidateForLogin())
		ctx.JSON(res.StatusCode, res)
		return
	}

	cashierObj.Id = cashierId
	if err := authCtrl.authSvc.CashierLogout(ctx, cashierObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}
