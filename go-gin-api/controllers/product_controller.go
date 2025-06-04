package controllers

import (
	"net/http"
	"time"
	"fmt"
	"go-gin-api/config"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
	"encoding/json"
)

// // Tüm ürünleri getir (ilişkili verilerle birlikte, Redis cache uyumlu)
// func GetProducts(c *gin.Context) {
// 	cacheKey := "products"

// 	// 1. Redis’te varsa onu getir
// 	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
// 	if err == nil {
// 		var products []models.Product
// 		if err := json.Unmarshal([]byte(cachedData), &products); err == nil {
// 			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
// 			c.JSON(http.StatusOK, products)
// 			return
// 		}
// 	}

// 	// 2. Veritabanından veri çek (ilişkili verilerle birlikte)
// 	var products []models.Product
// 	if err := config.DB.
// 		Preload("Origin").
// 		Preload("Location").
// 		Preload("SalesPoint").
// 		Find(&products).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanı hatası: " + err.Error()})
// 		return
// 	}

// 	// 3. Redis’e yaz (10 dakika süreli)
// 	jsonData, err := json.Marshal(products)
// 	if err == nil {
// 		if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
// 			fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
// 		} else {
// 			fmt.Println("⚠️ Redis yazma hatası:", err)
// 		}
// 	} else {
// 		fmt.Println("⚠️ JSON dönüşüm hatası:", err)
// 	}

// 	// 4. Yanıt döndür
// 	c.JSON(http.StatusOK, products)
// }



// Tüm ürünleri getir (ilişkili verilerle birlikte)
func GetProducts(c *gin.Context) {
	var products []models.Product
	config.DB.Preload("Origin").Preload("Location").Preload("SalesPoint").Find(&products)
	c.JSON(http.StatusOK, products)
}

// Yeni ürün oluştur (Redis cache uyumlu)
func CreateProduct(c *gin.Context) {
	var product models.Product

	// 1. JSON verisini parse et
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Tarih alanlarını kontrol et
	if product.ProductionDate.IsZero() {
		product.ProductionDate = time.Now()
	}
	if product.ExpirationDate.IsZero() {
		product.ExpirationDate = product.ProductionDate.AddDate(1, 0, 0)
	}

	// 3. Veritabanına kaydet
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün oluşturulamadı"})
		return
	}

	// 4. Redis cache temizlenir
	if err := config.RedisClient.Del(config.Ctx, "products").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'products' silindi.")
	} else {
		fmt.Println("⚠️ Redis cache silme hatası:", err)
	}

	// 5. Oluşturulan ürünü döndür
	c.JSON(http.StatusCreated, product)
}


// // Yeni ürün oluştur
// func CreateProduct(c *gin.Context) {
// 	var product models.Product
// 	if err := c.ShouldBindJSON(&product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Tarih alanlarının kontrolü
// 	if product.ProductionDate.IsZero() {
// 		product.ProductionDate = time.Now()
// 	}
// 	if product.ExpirationDate.IsZero() {
// 		product.ExpirationDate = product.ProductionDate.AddDate(1, 0, 0)
// 	}

// 	if err := config.DB.Create(&product).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün oluşturulamadı"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, product)
// }

// Tek ürün getir (Redis cache uyumlu)
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "product:" + id

	// 1. Redis'te varsa oradan getir
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var product models.Product
		if err := json.Unmarshal([]byte(cachedData), &product); err == nil {
			fmt.Println("📦 Redis cache'den getirildi:", cacheKey)
			c.JSON(http.StatusOK, product)
			return
		}
	}

	// 2. Veritabanından getir
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	// 3. Redis'e yaz
	jsonData, _ := json.Marshal(product)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
		fmt.Println("📝 Redis cache'e yazıldı:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis yazım hatası:", err)
	}

	// 4. Yanıt döndür
	c.JSON(http.StatusOK, product)
}


// // Tek ürün getir
// func GetProduct(c *gin.Context) {
// 	id := c.Param("id")
// 	var product models.Product
// 	if err := config.DB.First(&product, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, product)
// }

// Ürün güncelle (Redis cache uyumlu)
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "product:" + id

	var product models.Product

	// 1. Ürünü veritabanından al
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	// 2. Gelen JSON verisini işle
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Veritabanında güncelle
	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	// 4. Tekil ürün cache'ini güncelle
	jsonData, _ := json.Marshal(product)
	if err := config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err(); err == nil {
		fmt.Println("🔄 Redis cache güncellendi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis güncelleme hatası:", err)
	}

	// 5. Ürün listesi cache'ini sil
	if err := config.RedisClient.Del(config.Ctx, "products").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'products' silindi.")
	} else {
		fmt.Println("⚠️ Redis list cache silme hatası:", err)
	}

	// 6. Yanıt döndür
	c.JSON(http.StatusOK, product)
}


// // Ürün güncelle
// func UpdateProduct(c *gin.Context) {
// 	id := c.Param("id")
// 	var product models.Product
// 	if err := config.DB.First(&product, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	config.DB.Save(&product)
// 	c.JSON(http.StatusOK, product)
// }


// Ürün sil (Redis cache uyumlu)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "product:" + id

	var product models.Product

	// 1. Ürünün varlığını kontrol et
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	// 2. Veritabanından sil
	if err := config.DB.Delete(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün silinemedi"})
		return
	}

	// 3. Redis: tekil ürün cache'ini sil
	if err := config.RedisClient.Del(config.Ctx, cacheKey).Err(); err == nil {
		fmt.Println("🗑️ Redis cache silindi:", cacheKey)
	} else {
		fmt.Println("⚠️ Redis tekil silme hatası:", err)
	}

	// 4. Redis: ürün listesi cache'ini sil
	if err := config.RedisClient.Del(config.Ctx, "products").Err(); err == nil {
		fmt.Println("🧹 Redis cache 'products' silindi.")
	} else {
		fmt.Println("⚠️ Redis liste silme hatası:", err)
	}

	// 5. Yanıt döndür
	c.JSON(http.StatusOK, gin.H{"message": "Ürün silindi"})
}


// // Ürün sil
// func DeleteProduct(c *gin.Context) {
// 	id := c.Param("id")
// 	var product models.Product
// 	if err := config.DB.Delete(&product, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün silinemedi"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Ürün silindi"})
// }
