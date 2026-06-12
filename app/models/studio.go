package models

import (
	"github.com/goravel/framework/database/orm"
)

type Studio struct {
	orm.Model
	NamaStudio     string      `gorm:"column:nama_studio" json:"nama_studio"`
	Lokasi         string      `gorm:"column:lokasi" json:"lokasi"`
	Luas           int         `gorm:"column:luas" json:"luas"`
	JamOperasional *string     `gorm:"column:jam_operasional" json:"jam_operasional"` // Stored as JSON string
	HargaPerJam    int         `gorm:"column:harga_per_jam;default:75000" json:"harga_per_jam"`
	IsActive       bool        `gorm:"column:is_active;default:true" json:"is_active"`

	// Relationships
	AlatMusiks     []AlatMusik `gorm:"foreignKey:StudioID" json:"alat_musiks,omitempty"`
	Bookings       []Booking   `gorm:"foreignKey:StudioID" json:"bookings,omitempty"`
}

func (s *Studio) TableName() string {
	return "studios"
}
