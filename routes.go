package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/httpcontroller"
	"github.com/tamboto2000/gotoko-pos/httpmiddleware"
	"go.uber.org/zap"
)

func routes(log *zap.Logger) *gin.Engine {
	root := gin.New()

	// /cashiers
	cashiresCtrl := httpcontroller.NewCashierController(cashierSvc)
	cashiersGroup := root.Group("/cashiers")
	cashiersGroup.POST("", cashiresCtrl.Create)
	cashiersGroup.GET("/:cashierId", cashiresCtrl.GetDetail)
	cashiersGroup.GET("", cashiresCtrl.GetList)
	cashiersGroup.PUT("/:cashierId", cashiresCtrl.Update)

	authCtrl := httpcontroller.NewAuthController(authSvc)
	authGroup := root.Group("/cashiers")
	authGroup.GET("/:cashierId/passcode", authCtrl.GetCashierPasscode)
	authGroup.POST("/:cashierId/login", authCtrl.CashierLogin)
	authGroup.POST("/:cashierId/logout", authCtrl.CashierLogout)

	// /products group without authorization
	prodCtrl := httpcontroller.NewProductsController(prodSvc)
	productGroup := root.Group("/products")
	// productGroup.Use(httpmiddleware.AuthMiddleware(cashiersRepo, cashiersMRepo, log))
	productGroup.POST("", prodCtrl.CreateProduct)
	productGroup.PUT("/:productId", prodCtrl.UpdateProduct)
	productGroup.DELETE("/:productId", prodCtrl.DeleteProduct)
	productGroup.GET("/test", func(ctx *gin.Context) {
		ctx.Data(200, "application/text", []byte("/product route"))
	})

	// /products group with authorization
	productGroupAuth := root.Group("/products")
	productGroupAuth.Use(httpmiddleware.AuthMiddleware(cashiersRepo, cashiersMRepo, log))
	productGroupAuth.GET("/:productId", prodCtrl.GetProductDetail)
	productGroupAuth.GET("", prodCtrl.GetProductList)

	// /categories group without authorization
	categoryGroup := root.Group("/categories")
	categoryGroup.POST("", prodCtrl.CreateCategory)

	// /categories group with authorization
	categoryGroupAuth := root.Group("/categories")
	categoryGroupAuth.Use(httpmiddleware.AuthMiddleware(cashiersRepo, cashiersMRepo, log))
	categoryGroupAuth.GET("/:categoryId", prodCtrl.GetCategoryDetail)
	categoryGroupAuth.GET("", prodCtrl.GetCategoryList)

	return root
}
