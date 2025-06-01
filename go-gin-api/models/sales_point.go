package models

type SalesPoint struct {
	SalesPointID uint   `gorm:"primaryKey" json:"sales_point_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
}
