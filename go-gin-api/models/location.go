package models

type Location struct {
	LocationID uint   `gorm:"primaryKey" json:"location_id"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
}
