package models

import (
	"github.com/goravel/framework/database/orm"
)

type Customer struct {
	orm.Model
	UserID   uint      `gorm:"column:user_id" json:"user_id"`
	Nama     string    `gorm:"column:nama" json:"nama"`
	Telepon  *string   `gorm:"column:telepon" json:"telepon"`

	// Relationships
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Bookings []Booking `gorm:"foreignKey:CustomerID" json:"bookings,omitempty"`
}

func (c *Customer) TableName() string {
	return "customers"
}
