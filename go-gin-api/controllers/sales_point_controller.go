package controllers

import (
	"net/http"
	"go-gin-api/config"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
)

// Tüm satış noktalarını getir
func GetSalesPoints(c *gin.Context) {
	var salesPoints []models.SalesPoint
	config.DB.Find(&salesPoints)
	c.JSON(http.StatusOK, salesPoints)
}

// Yeni satış noktası oluştur
func CreateSalesPoint(c *gin.Context) {
	var salesPoint models.SalesPoint
	if err := c.ShouldBindJSON(&salesPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&salesPoint)
	c.JSON(http.StatusCreated, salesPoint)
}

// Tek bir satış noktasını getir
func GetSalesPoint(c *gin.Context) {
	id := c.Param("id")
	var salesPoint models.SalesPoint
	if err := config.DB.First(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, salesPoint)
}

// Satış noktasını güncelle
func UpdateSalesPoint(c *gin.Context) {
	id := c.Param("id")
	var salesPoint models.SalesPoint
	if err := config.DB.First(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
		return
	}
	if err := c.ShouldBindJSON(&salesPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&salesPoint)
	c.JSON(http.StatusOK, salesPoint)
}

// Satış noktasını sil
func DeleteSalesPoint(c *gin.Context) {
	id := c.Param("id")
	var salesPoint models.SalesPoint
	if err := config.DB.Delete(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinemedi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Satış noktası silindi"})
}
