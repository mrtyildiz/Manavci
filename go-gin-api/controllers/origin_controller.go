package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"go-gin-api/config"
	"go-gin-api/models"
	"fmt"
	"encoding/json"
	"time"
)

// Tüm origin verilerini getir (Redis cache uyumlu)
func GetOrigins(c *gin.Context) {
	cacheKey := "origins"

	// 1. Redis cache kontrolü
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var origins []models.Origin
		if err := json.Unmarshal([]byte(cachedData), &origins); err == nil {
			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
			c.JSON(http.StatusOK, origins)
			return
		}
	}

	// 2. Veritabanından veri çek
	var origins []models.Origin
	if err := config.DB.Find(&origins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanı hatası"})
		return
	}

	// 3. Redis'e yaz
	jsonData, _ := json.Marshal(origins)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
		fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis yazım hatası:", err)
	}

	// 4. Yanıtı döndür
	c.JSON(http.StatusOK, origins)
}


// func GetOrigins(c *gin.Context) {
// 	var origins []models.Origin
// 	config.DB.Find(&origins)
// 	c.JSON(http.StatusOK, origins)
// }

// Yeni origin oluştur (Redis cache uyumlu)
func CreateOrigin(c *gin.Context) {
	var origin models.Origin

	// 1. Gelen JSON verisini parse et
	if err := c.ShouldBindJSON(&origin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Veritabanına kaydet
	if err := config.DB.Create(&origin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanı hatası"})
		return
	}

	// 3. Redis'teki 'origins' cache’ini sil
	if err := config.RedisClient.Del(config.Ctx, "origins").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'origins' silindi.")
	} else {
		fmt.Println("⚠️ Redis cache silinirken hata:", err)
	}

	// 4. Yanıtı döndür
	c.JSON(http.StatusCreated, origin)
}


// func CreateOrigin(c *gin.Context) {
// 	var origin models.Origin
// 	if err := c.ShouldBindJSON(&origin); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	config.DB.Create(&origin)
// 	c.JSON(http.StatusCreated, origin)
// }

// Tek bir origin getir (Redis cache uyumlu)
func GetOrigin(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "origin:" + id

	// 1. Redis'te varsa onu getir
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var origin models.Origin
		if err := json.Unmarshal([]byte(cachedData), &origin); err == nil {
			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
			c.JSON(http.StatusOK, origin)
			return
		}
	}

	// 2. Veritabanından çek
	var origin models.Origin
	if err := config.DB.First(&origin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
		return
	}

	// 3. Redis'e yaz (10 dakikalığına)
	jsonData, _ := json.Marshal(origin)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
		fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis yazım hatası:", err)
	}

	// 4. Yanıt döndür
	c.JSON(http.StatusOK, origin)
}


// func GetOrigin(c *gin.Context) {
// 	id := c.Param("id")
// 	var origin models.Origin
// 	if err := config.DB.First(&origin, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, origin)
// }

// Origin güncelle (Redis cache uyumlu)
func UpdateOrigin(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "origin:" + id

	var origin models.Origin

	// 1. Mevcut veriyi veritabanından al
	if err := config.DB.First(&origin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
		return
	}

	// 2. Gelen JSON verisi ile eşleştir
	if err := c.ShouldBindJSON(&origin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Veritabanında güncelle
	if err := config.DB.Save(&origin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	// 4. Redis'teki 'origin:{id}' güncellenir
	originJSON, _ := json.Marshal(origin)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, originJSON, 10*time.Minute).Err(); err == nil {
		fmt.Println("🔄 Redis cache güncellendi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis güncelleme hatası:", err)
	}

	// 5. 'origins' listesi cache'ini de sil
	if err := config.RedisClient.Del(config.Ctx, "origins").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'origins' silindi.")
	} else {
		fmt.Println("⚠️ Redis silme hatası:", err)
	}

	// 6. Güncellenmiş veriyi döndür
	c.JSON(http.StatusOK, origin)
}


// func UpdateOrigin(c *gin.Context) {
// 	id := c.Param("id")
// 	var origin models.Origin
// 	if err := config.DB.First(&origin, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
// 		return
// 	}
// 	if err := c.ShouldBindJSON(&origin); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	config.DB.Save(&origin)
// 	c.JSON(http.StatusOK, origin)
// }

// Origin sil (Redis cache uyumlu)
func DeleteOrigin(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "origin:" + id

	var origin models.Origin

	// 1. Önce verinin var olup olmadığını kontrol et
	if err := config.DB.First(&origin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Origin bulunamadı"})
		return
	}

	// 2. Veritabanından sil
	if err := config.DB.Delete(&origin, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"})
		return
	}

	// 3. Redis: Tekil cache'i sil
	if err := config.RedisClient.Del(config.Ctx, cacheKey).Err(); err == nil {
		fmt.Println("🗑️ Redis cache silindi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis cache silme hatası:", err)
	}

	// 4. Redis: Liste cache'i sil
	if err := config.RedisClient.Del(config.Ctx, "origins").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'origins' silindi.")
	} else {
		fmt.Println("⚠️ Redis list cache silme hatası:", err)
	}

	// 5. Yanıt döndür
	c.JSON(http.StatusOK, gin.H{"message": "Origin silindi"})
}


// func DeleteOrigin(c *gin.Context) {
// 	id := c.Param("id")
// 	var origin models.Origin
// 	if err := config.DB.Delete(&origin, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Silinemedi"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Silindi"})
// }
