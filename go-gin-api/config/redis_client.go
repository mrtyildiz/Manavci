package config

import (
	"context"
	"log"
	"os"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR") // .env ya da docker-compose üzerinden gelir
	if addr == "" {
		addr = "redis:6379" // varsayılan
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // parola varsa buraya ekle
		DB:       0,  // default DB
	})

	// Bağlantıyı test et
	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		log.Fatalf("❌ Redis bağlantı hatası: %v", err)
	} else {
		log.Println("✅ Redis bağlantısı başarılı.")
	}
}
