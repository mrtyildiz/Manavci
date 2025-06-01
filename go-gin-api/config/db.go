package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-gin-api/models" // tüm modelleri burada import et
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanı bağlantısı başarısız:", err)
	}

	// 🔽 Buraya AutoMigrate eklenir
	err = DB.AutoMigrate(
		&models.Origin{},
		&models.Location{},
		&models.SalesPoint{},
		&models.Product{},
	)

	if err != nil {
		log.Fatal("Migration sırasında hata oluştu:", err)
	}

	fmt.Println("Veritabanına başarıyla bağlandı ve tablolar migrate edildi!")
}
