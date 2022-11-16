package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type OrderDetails struct {
	ID        uint64                `gorm:"primary_key:auto_increment" json:"id"`
	OrderId   uint64                `gorm:"type:integer" json:"order_id"`
	Item      string                `gorm:"type:varchar(255)" json:"item"`
	Qty       int64                 `gorm:"type:integer" json:"qty"`
	Price     float64               `gorm:"type:float" json:"price"`
	Amount    float64               `gorm:"type:float" json:"amount"`
	CreatedAt time.Time             `gorm:"<-:create"`
	UpdatedAt time.Time             `gorm:"<-:update"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
