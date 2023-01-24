package httpcontroller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/domain/products"
	"github.com/tamboto2000/gotoko-pos/helpers/httpresponse"
	"github.com/tamboto2000/gotoko-pos/services/product"
)

type ProductsController struct {
	prodSvc *product.ProductService
}

func NewProductsController(prodSvc *product.ProductService) *ProductsController {
	return &ProductsController{prodSvc: prodSvc}
}

func (prodCtrl *ProductsController) CreateProduct(ctx *gin.Context) {
	prodObj := new(products.Product)
	if err := ctx.BindJSON(prodObj); err != nil {
		res := httpresponse.FromError(prodObj.Validate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	if err := prodCtrl.prodSvc.CreateProduct(ctx, prodObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(prodObj)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) GetProductDetail(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	productIdStr := ctx.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		errl.Add(apperror.New(`"productId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "productId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	prod, err := prodCtrl.prodSvc.GetProductDetail(ctx, productId)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(prod)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) GetProductList(ctx *gin.Context) {
	limit := 0
	skip := 0
	categoryId := 0
	q := ctx.Query("q")

	// validate inputs
	errl := apperror.NewErrorList()
	errl.SetPrefix("query ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	if s := ctx.Query("limit"); s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			errl.Add(apperror.New(`"limit" is invalid, valid input are positive integer`, apperror.AnyInvalid, "limit"))
		} else {
			limit = i
		}
	}

	if s := ctx.Query("skip"); s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			errl.Add(apperror.New(`"skip" is invalid, valid input are positive integer`, apperror.AnyInvalid, "skip"))
		} else {
			skip = i
		}
	}

	if s := ctx.Query("categoryId"); s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			errl.Add(apperror.New(`"categoryId" is invalid, valid input are positive integer`, apperror.AnyInvalid, "categoryId"))
		} else {
			categoryId = i
		}
	}

	if err := errl.BuildError(); err != nil {
		log.Println("ERROR validation: ", err.Error())
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	prodList, err := prodCtrl.prodSvc.GetProductList(ctx, limit, skip, categoryId, q)
	if err != nil {
		log.Println("ERROR prodSvc.GetProductList: ", err.Error())
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(prodList)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) UpdateProduct(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	productIdStr := ctx.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		errl.Add(apperror.New(`"productId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "productId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	prodObj := new(products.Product)
	if err := ctx.BindJSON(prodObj); err != nil {
		res := httpresponse.FromError(prodObj.ValidateForUpdate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	prodObj.Id = productId
	if err := prodCtrl.prodSvc.UpdateProduct(ctx, prodObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) DeleteProduct(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	productIdStr := ctx.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		errl.Add(apperror.New(`"productId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "productId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	if err := prodCtrl.prodSvc.DeleteProduct(ctx, productId); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) CreateCategory(ctx *gin.Context) {
	catObj := new(products.Category)
	if err := ctx.BindJSON(catObj); err != nil {
		res := httpresponse.FromError(catObj.Validate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	if err := prodCtrl.prodSvc.CreateCategory(ctx, catObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(catObj)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) GetCategoryDetail(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	categoryIdStr := ctx.Param("categoryId")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		errl.Add(apperror.New(`"categoryId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "categoryId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	cat, err := prodCtrl.prodSvc.GetCategoryDetail(ctx, categoryId)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(cat)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) GetCategoryList(ctx *gin.Context) {
	limit := 0
	skip := 0

	// validate inputs
	errl := apperror.NewErrorList()
	errl.SetPrefix("query ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	if s := ctx.Query("limit"); s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			errl.Add(apperror.New(`"limit" is invalid, valid input are positive integer`, apperror.AnyInvalid, "limit"))
		} else {
			limit = i
		}
	}

	if s := ctx.Query("skip"); s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			errl.Add(apperror.New(`"skip" is invalid, valid input are positive integer`, apperror.AnyInvalid, "skip"))
		} else {
			skip = i
		}
	}

	catList, err := prodCtrl.prodSvc.GetCategoryList(ctx, limit, skip)
	if err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(catList)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) UpdateCategory(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	categoryIdStr := ctx.Param("categoryId")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		errl.Add(apperror.New(`"categoryId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "categoryId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	catObj := new(products.Category)

	if err := ctx.BindJSON(catObj); err != nil {
		res := httpresponse.FromError(catObj.ValidateForUpdate())
		ctx.JSON(res.StatusCode, res)

		return
	}

	catObj.Id = categoryId
	if err := prodCtrl.prodSvc.UpdateCategory(ctx, catObj); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}

func (prodCtrl *ProductsController) DeleteCategory(ctx *gin.Context) {
	errl := apperror.NewErrorList()
	errl.SetPrefix("param ValidationError: ")
	errl.SetType(apperror.AnyInvalid)

	categoryIdStr := ctx.Param("categoryId")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		errl.Add(apperror.New(`"categoryId" is invalid, valid input is number (0-9)`, apperror.AnyInvalid, "categoryId"))
	}

	if err := errl.BuildError(); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	if err := prodCtrl.prodSvc.DeleteCategory(ctx, categoryId); err != nil {
		res := httpresponse.FromError(err)
		ctx.JSON(res.StatusCode, res)

		return
	}

	res := httpresponse.Success(nil)
	ctx.JSON(res.StatusCode, res)
}
