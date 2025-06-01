package models

type Origin struct {
	OriginID    uint   `gorm:"primaryKey" json:"origin_id"`
	OriginName  string `json:"origin_name"`
	Description string `json:"description"`
}

