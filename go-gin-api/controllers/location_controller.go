package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"go-gin-api/config"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
	"fmt"
)

// Tüm lokasyonları getir (Redis cache destekli)
func GetLocations(c *gin.Context) {
	cacheKey := "locations"

	// 1. Redis cache kontrolü
	cached, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		// Cache bulunduysa direkt onu döndür
		var cachedLocations []models.Location
		//fmt.Println("📌 fmt: /test endpoint çağrıldı")
		if err := json.Unmarshal([]byte(cached), &cachedLocations); err == nil {
			c.JSON(http.StatusOK, cachedLocations)
			return
		}
	}

	// 2. PostgreSQL'den veriyi çek
	var locations []models.Location
	result := config.DB.Find(&locations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanı hatası"})
		return
	}
	//fmt.Println("Redis ileyim ya sen") 
	// 3. JSON'a çevirip Redis'e yaz (5 dk süreyle)
	jsonData, _ := json.Marshal(locations)
	config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 5*time.Minute)

	// 4. Yanıt olarak döndür
	c.JSON(http.StatusOK, locations)
}

// // Tüm lokasyonları getir
// func GetLocations_DB(c *gin.Context) {
// 	var locations []models.Location
// 	config.DB.Find(&locations)
// 	c.JSON(http.StatusOK, locations)
// }

// Yeni lokasyon oluştur (Redis cache uyumlu)
func CreateLocation(c *gin.Context) {
	var location models.Location

	// JSON verisini parse et
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Veritabanına kaydet
	if err := config.DB.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanı hatası"})
		return
	}

	// Redis cache temizlenir
	cacheKey := "locations"
	if err := config.RedisClient.Del(config.Ctx, cacheKey).Err(); err == nil {
		fmt.Println("🧹 Redis cache 'locations' silindi.")
	} else {
		fmt.Println("⚠️ Redis cache silinirken hata:", err)
	}

	// Yanıtı döndür
	c.JSON(http.StatusCreated, location)
}


// // Yeni lokasyon oluştur
// func CreateLocation(c *gin.Context) {
// 	var location models.Location
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	config.DB.Create(&location)
// 	c.JSON(http.StatusCreated, location)
// }

// Tek bir lokasyonu getir (Redis cache uyumlu)
func GetLocation(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "location:" + id

	// Redis cache kontrolü
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var location models.Location
		if err := json.Unmarshal([]byte(cachedData), &location); err == nil {
			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
			c.JSON(http.StatusOK, location)
			return
		}
	}

	// Veritabanından getir
	var location models.Location
	if err := config.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
		return
	}

	// Redis'e kaydet (10 dakika süreyle)
	locationJSON, _ := json.Marshal(location)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, locationJSON, 10*time.Minute).Err(); err == nil {
		fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis yazım hatası:", err)
	}

	// Yanıtı döndür
	c.JSON(http.StatusOK, location)
}


// // Tek bir lokasyonu getir
// func GetLocation(c *gin.Context) {
// 	id := c.Param("id")
// 	var location models.Location
// 	if err := config.DB.First(&location, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, location)
// }

// Lokasyonu güncelle (Redis cache uyumlu)
func UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "location:" + id

	var location models.Location

	// Veritabanından mevcut veriyi al
	if err := config.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
		return
	}

	// Gelen JSON ile eşleştir
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Veriyi güncelle
	if err := config.DB.Save(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	// Redis cache güncellenir
	locationJSON, _ := json.Marshal(location)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, locationJSON, 10*time.Minute).Err(); err == nil {
		fmt.Println("🔄 Redis cache güncellendi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis güncelleme hatası:", err)
	}

	// Yanıtı döndür
	c.JSON(http.StatusOK, location)
}


// // Lokasyonu güncelle
// func UpdateLocation(c *gin.Context) {
// 	id := c.Param("id")
// 	var location models.Location
// 	if err := config.DB.First(&location, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
// 		return
// 	}
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	config.DB.Save(&location)
// 	c.JSON(http.StatusOK, location)
// }

// Lokasyonu sil (Redis cache uyumlu)
func DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "location:" + id

	var location models.Location

	// Önce lokasyonun varlığını kontrol et
	if err := config.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasyon bulunamadı"})
		return
	}

	// Veritabanından sil
	if err := config.DB.Delete(&location, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"})
		return
	}

	// Redis cache sil
	if err := config.RedisClient.Del(config.Ctx, cacheKey).Err(); err == nil {
		fmt.Println("🗑️ Redis cache silindi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis silme hatası:", err)
	}

	// Yanıt döndür
	c.JSON(http.StatusOK, gin.H{"message": "Lokasyon silindi"})
}


// // Lokasyonu sil
// func DeleteLocation(c *gin.Context) {
// 	id := c.Param("id")
// 	var location models.Location
// 	if err := config.DB.Delete(&location, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Silinemedi"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Lokasyon silindi"})
// }
