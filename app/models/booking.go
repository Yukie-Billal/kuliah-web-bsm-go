package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Booking struct {
	orm.Model
	BookingID   string     `gorm:"column:booking_id;unique" json:"booking_id"`
	CustomerID  uint       `gorm:"column:customer_id" json:"customer_id"`
	StudioID    uint       `gorm:"column:studio_id" json:"studio_id"`
	Tanggal     time.Time  `gorm:"column:tanggal;type:date" json:"tanggal"`
	Jam         string     `gorm:"column:jam" json:"jam"`
	Durasi      int        `gorm:"column:durasi;default:1" json:"durasi"`
	TotalBiaya  int        `gorm:"column:total_biaya;default:75000" json:"total_biaya"`
	Status      string     `gorm:"column:status;default:pending" json:"status"`
	Catatan     *string    `gorm:"column:catatan" json:"catatan"`
	CheckedInAt *time.Time `gorm:"column:checked_in_at" json:"checked_in_at"`
	ApprovedBy  *uint      `gorm:"column:approved_by" json:"approved_by"`

	// Relationships
	Customer    *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Studio      *Studio    `gorm:"foreignKey:StudioID" json:"studio,omitempty"`
	Approver    *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

func (b *Booking) TableName() string {
	return "bookings"
}
