package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Customer struct {
	ID        uint64                `gorm:"primary_key:auto_increment" json:"id"`
	Name      string                `gorm:"type:varchar(255)" json:"name"`
	Email     string                `gorm:"unique:varchar(255)" json:"email"`
	Password  string                `gorm:"->;<-;not null " json:"-"`
	CreatedAt time.Time             `gorm:"<-:create"`
	UpdatedAt time.Time             `gorm:"<-:update"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
