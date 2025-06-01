package controllers

import (
	"net/http"
	"go-gin-api/config"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
)

// Tüm lokasyonları getir
func GetLocations(c *gin.Context) {
	var locations []models.Location
	config.DB.Find(&locations)
	c.JSON(http.StatusOK, locations)
}

// Yeni lokasyon oluştur
func CreateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&location)
	c.JSON(http.StatusCreated, location)
}

// Tek bir lokasyonu getir
func GetLocation(c *gin.Context) {
	id := c.Param("id")
	var location models.Location
	if err := config.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, location)
}

// Lokasyonu güncelle
func UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	var location models.Location
	if err := config.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
		return
	}
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&location)
	c.JSON(http.StatusOK, location)
}

// Lokasyonu sil
func DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	var location models.Location
	if err := config.DB.Delete(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinemedi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lokasyon silindi"})
}
