package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"go-gin-api/config"
	"go-gin-api/models"
)

func GetOrigins(c *gin.Context) {
	var origins []models.Origin
	config.DB.Find(&origins)
	c.JSON(http.StatusOK, origins)
}

func CreateOrigin(c *gin.Context) {
	var origin models.Origin
	if err := c.ShouldBindJSON(&origin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&origin)
	c.JSON(http.StatusCreated, origin)
}

func GetOrigin(c *gin.Context) {
	id := c.Param("id")
	var origin models.Origin
	if err := config.DB.First(&origin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, origin)
}

func UpdateOrigin(c *gin.Context) {
	id := c.Param("id")
	var origin models.Origin
	if err := config.DB.First(&origin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
		return
	}
	if err := c.ShouldBindJSON(&origin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&origin)
	c.JSON(http.StatusOK, origin)
}

func DeleteOrigin(c *gin.Context) {
	id := c.Param("id")
	var origin models.Origin
	if err := config.DB.Delete(&origin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinemedi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Silindi"})
}
