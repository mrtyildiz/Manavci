package models

import "time"

type Product struct {
	ProductID       uint      `gorm:"primaryKey" json:"product_id"`
	ProductName     string    `json:"product_name"`
	Price           float64   `json:"price"`
	Stock           int       `json:"stock"`
	ProductionDate  time.Time `json:"production_date"`
	ExpirationDate  time.Time `json:"expiration_date"`

	OriginID         uint `json:"origin_id"`
	CurrentLocationID uint `json:"current_location_id"`
	SalesPointID     uint `json:"sales_point_id"`

	// Origin           Origin `gorm:"foreignKey:OriginID" json:"origin,omitempty"`
	// Location   Location   `gorm:"foreignKey:CurrentLocationID" json:"location,omitempty"`
	// SalesPoint SalesPoint `gorm:"foreignKey:SalesPointID" json:"sales_point,omitempty"`

}