package controllers

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"go-gin-api/config"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
)

// Tüm satış noktalarını getir (Redis cache uyumlu)
func GetSalesPoints(c *gin.Context) {
	cacheKey := "salespoints"

	// 1. Redis cache kontrolü
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var salesPoints []models.SalesPoint
		if err := json.Unmarshal([]byte(cachedData), &salesPoints); err == nil {
			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
			c.JSON(http.StatusOK, salesPoints)
			return
		}
	}

	// 2. Veritabanından veriyi çek
	var salesPoints []models.SalesPoint
	if err := config.DB.Find(&salesPoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanı hatası: " + err.Error()})
		return
	}

	// 3. Redis'e yaz (10 dakikalığına)
	jsonData, err := json.Marshal(salesPoints)
	if err == nil {
		if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
			fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
		} else {
			fmt.Println("⚠️ Redis yazım hatası:", err)
		}
	} else {
		fmt.Println("⚠️ JSON dönüşüm hatası:", err)
	}

	// 4. Yanıtı döndür
	c.JSON(http.StatusOK, salesPoints)
}


// // Tüm satış noktalarını getir
// func GetSalesPoints(c *gin.Context) {
// 	var salesPoints []models.SalesPoint
// 	config.DB.Find(&salesPoints)
// 	c.JSON(http.StatusOK, salesPoints)
// }

// Yeni satış noktası oluştur (Redis cache uyumlu)
func CreateSalesPoint(c *gin.Context) {
	var salesPoint models.SalesPoint

	// 1. JSON'dan parse et
	if err := c.ShouldBindJSON(&salesPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Veritabanına kaydet
	if err := config.DB.Create(&salesPoint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Satış noktası oluşturulamadı"})
		return
	}

	// 3. Redis cache sil
	if err := config.RedisClient.Del(config.Ctx, "salespoints").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'salespoints' silindi.")
	} else {
		fmt.Println("⚠️ Redis silme hatası:", err)
	}

	// 4. Oluşturulan veriyi döndür
	c.JSON(http.StatusCreated, salesPoint)
}

// // Yeni satış noktası oluştur
// func CreateSalesPoint(c *gin.Context) {
// 	var salesPoint models.SalesPoint
// 	if err := c.ShouldBindJSON(&salesPoint); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	config.DB.Create(&salesPoint)
// 	c.JSON(http.StatusCreated, salesPoint)
// }

// Tek bir satış noktasını getir (Redis cache uyumlu)
func GetSalesPoint(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "salespoint:" + id

	// 1. Redis cache kontrolü
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var salesPoint models.SalesPoint
		if err := json.Unmarshal([]byte(cachedData), &salesPoint); err == nil {
			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
			c.JSON(http.StatusOK, salesPoint)
			return
		}
	}

	// 2. Veritabanından getir
	var salesPoint models.SalesPoint
	if err := config.DB.First(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
		return
	}

	// 3. Redis'e yaz
	jsonData, err := json.Marshal(salesPoint)
	if err == nil {
		if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
			fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
		} else {
			fmt.Println("⚠️ Redis yazım hatası:", err)
		}
	} else {
		fmt.Println("⚠️ JSON dönüşüm hatası:", err)
	}

	// 4. Yanıt döndür
	c.JSON(http.StatusOK, salesPoint)
}

// // Tek bir satış noktasını getir
// func GetSalesPoint(c *gin.Context) {
// 	id := c.Param("id")
// 	var salesPoint models.SalesPoint
// 	if err := config.DB.First(&salesPoint, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, salesPoint)
// }

// Satış noktasını güncelle (Redis cache uyumlu)
func UpdateSalesPoint(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "salespoint:" + id

	var salesPoint models.SalesPoint

	// 1. Veritabanından mevcut kaydı al
	if err := config.DB.First(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
		return
	}

	// 2. Gelen veriyi bind et
	if err := c.ShouldBindJSON(&salesPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Veritabanında güncelle
	if err := config.DB.Save(&salesPoint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	// 4. Redis: tekil cache güncelle
	jsonData, _ := json.Marshal(salesPoint)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
		fmt.Println("🔄 Redis cache güncellendi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis güncelleme hatası:", err)
	}

	// 5. Redis: liste cache silinsin
	if err := config.RedisClient.Del(config.Ctx, "salespoints").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'salespoints' silindi.")
	} else {
		fmt.Println("⚠️ Redis liste silme hatası:", err)
	}

	// 6. Yanıt döndür
	c.JSON(http.StatusOK, salesPoint)
}


// // Satış noktasını güncelle
// func UpdateSalesPoint(c *gin.Context) {
// 	id := c.Param("id")
// 	var salesPoint models.SalesPoint
// 	if err := config.DB.First(&salesPoint, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
// 		return
// 	}
// 	if err := c.ShouldBindJSON(&salesPoint); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	config.DB.Save(&salesPoint)
// 	c.JSON(http.StatusOK, salesPoint)
// }

// Satış noktasını sil (Redis cache uyumlu)
func DeleteSalesPoint(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "salespoint:" + id

	var salesPoint models.SalesPoint

	// 1. Kayıt var mı kontrol et
	if err := config.DB.First(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satış noktası bulunamadı"})
		return
	}

	// 2. Veritabanından sil
	if err := config.DB.Delete(&salesPoint, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi"})
		return
	}

	// 3. Redis: tekil cache sil
	if err := config.RedisClient.Del(config.Ctx, cacheKey).Err(); err == nil {
		fmt.Println("🗑️ Redis cache silindi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis tekil silme hatası:", err)
	}

	// 4. Redis: liste cache sil
	if err := config.RedisClient.Del(config.Ctx, "salespoints").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'salespoints' silindi.")
	} else {
		fmt.Println("⚠️ Redis liste silme hatası:", err)
	}

	// 5. Yanıt döndür
	c.JSON(http.StatusOK, gin.H{"message": "Satış noktası silindi"})
}


// // Satış noktasını sil
// func DeleteSalesPoint(c *gin.Context) {
// 	id := c.Param("id")
// 	var salesPoint models.SalesPoint
// 	if err := config.DB.Delete(&salesPoint, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Silinemedi"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Satış noktası silindi"})
// }
