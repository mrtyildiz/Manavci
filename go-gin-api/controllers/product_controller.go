package controllers

import (
	"net/http"
	"time"

	"go-gin-api/config"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
)

// Tüm ürünleri getir (ilişkili verilerle birlikte)
func GetProducts(c *gin.Context) {
	var products []models.Product
	config.DB.Preload("Origin").Preload("Location").Preload("SalesPoint").Find(&products)
	c.JSON(http.StatusOK, products)
}

// Yeni ürün oluştur
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tarih alanlarının kontrolü
	if product.ProductionDate.IsZero() {
		product.ProductionDate = time.Now()
	}
	if product.ExpirationDate.IsZero() {
		product.ExpirationDate = product.ProductionDate.AddDate(1, 0, 0)
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün oluşturulamadı"})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// Tek ürün getir
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Ürün güncelle
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

// Ürün sil
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.Delete(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün silinemedi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ürün silindi"})
}
