package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/httpcontroller"
	"github.com/tamboto2000/gotoko-pos/httpmiddleware"
	"go.uber.org/zap"
)

func routes(log *zap.Logger) *gin.Engine {
	root := gin.New()
	root.Use(httpmiddleware.RequestLog(log))

	// /cashiers
	cashiresCtrl := httpcontroller.NewCashierController(cashierSvc)
	cashiersGroup := root.Group("/cashiers")
	cashiersGroup.POST("", cashiresCtrl.Create)
	cashiersGroup.GET("/:cashierId", cashiresCtrl.GetDetail)
	cashiersGroup.GET("", cashiresCtrl.GetList)
	cashiersGroup.PUT("/:cashierId", cashiresCtrl.Update)
	cashiersGroup.DELETE("/:cashierId", cashiresCtrl.Delete)

	authCtrl := httpcontroller.NewAuthController(authSvc)
	authGroup := root.Group("/cashiers")
	authGroup.GET("/:cashierId/passcode", authCtrl.GetCashierPasscode)
	authGroup.POST("/:cashierId/login", authCtrl.CashierLogin)
	authGroup.POST("/:cashierId/logout", authCtrl.CashierLogout)

	// /products group without authorization
	prodCtrl := httpcontroller.NewProductsController(prodSvc)
	productGroup := root.Group("/products")
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
	categoryGroup.PUT("/:categoryId", prodCtrl.UpdateCategory)
	categoryGroup.DELETE("/:categoryId", prodCtrl.DeleteCategory)

	// /categories group with authorization
	categoryGroupAuth := root.Group("/categories")
	categoryGroupAuth.Use(httpmiddleware.AuthMiddleware(cashiersRepo, cashiersMRepo, log))
	categoryGroupAuth.GET("/:categoryId", prodCtrl.GetCategoryDetail)
	categoryGroupAuth.GET("", prodCtrl.GetCategoryList)

	// /payments without authorization
	payCtrl := httpcontroller.NewPaymentController(paySvc)
	paymentGroup := root.Group("/payments")
	paymentGroup.POST("", payCtrl.CreatePayment)
	paymentGroup.PUT("/:paymentId", payCtrl.UpdatePayment)
	paymentGroup.DELETE("/:paymentId", payCtrl.DeletePayment)

	// /payments with authorization
	paymentGroupAuth := root.Group("/payments")
	paymentGroupAuth.Use(httpmiddleware.AuthMiddleware(cashiersRepo, cashiersMRepo, log))
	paymentGroupAuth.GET("/:paymentId", payCtrl.GetPaymentDetail)

	// /orders with authorization
	orderCtrl := httpcontroller.NewOrdersController(orderSvc)
	orderGroupAuth := root.Group("/orders")
	orderGroupAuth.Use(httpmiddleware.AuthMiddleware(cashiersRepo, cashiersMRepo, log))
	orderGroupAuth.POST("/subtotal", orderCtrl.GetOrderSubtotal)
	orderGroupAuth.POST("", orderCtrl.AddOrder)

	return root
}
