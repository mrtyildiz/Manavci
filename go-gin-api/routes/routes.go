package routes

import (
	"github.com/gin-gonic/gin"
	"go-gin-api/controllers"
)

func SetupRoutes(router *gin.Engine) {
	// Origin routes
	router.GET("/origins", controllers.GetOrigins)
	router.POST("/origins", controllers.CreateOrigin)
	router.GET("/origins/:id", controllers.GetOrigin)
	router.PUT("/origins/:id", controllers.UpdateOrigin)
	router.DELETE("/origins/:id", controllers.DeleteOrigin)

		// Location routes
	router.GET("/locations", controllers.GetLocations)
	router.POST("/locations", controllers.CreateLocation)
	router.GET("/locations/:id", controllers.GetLocation)
	router.PUT("/locations/:id", controllers.UpdateLocation)
	router.DELETE("/locations/:id", controllers.DeleteLocation)
	// Diğer tablolar için aynı yapıyı location, sales_point, product için uygula.
	// SalesPoint routes
	router.GET("/sales-points", controllers.GetSalesPoints)
	router.POST("/sales-points", controllers.CreateSalesPoint)
	router.GET("/sales-points/:id", controllers.GetSalesPoint)
	router.PUT("/sales-points/:id", controllers.UpdateSalesPoint)
	router.DELETE("/sales-points/:id", controllers.DeleteSalesPoint)
	// Products
	router.GET("/products", controllers.GetProducts)
	router.POST("/products", controllers.CreateProduct)
	router.GET("/products/:id", controllers.GetProduct)
	router.PUT("/products/:id", controllers.UpdateProduct)
	router.DELETE("/products/:id", controllers.DeleteProduct)
}
