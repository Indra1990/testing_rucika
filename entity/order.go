package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Orders struct {
	ID           uint64                `gorm:"primary_key:auto_increment" json:"id"`
	CustomerId   uint64                `gorm:"type:integer" json:"customerId"`
	Title        string                `gorm:"type:varchar(255)" json:"title"`
	OrderNumber  string                `gorm:"type:varchar(255)" json:"orderNumber"`
	Note         string                `gorm:"type:text" json:"note"`
	Total        float64               `gorm:"type:float" json:"total"`
	CreatedAt    time.Time             `gorm:"<-:create"`
	UpdatedAt    time.Time             `gorm:"<-:update"`
	DeletedAt    soft_delete.DeletedAt `gorm:"softDelete:flag"`
	OrderDetails []OrderDetails        `gorm:"foreignKey:OrderId"`
}
