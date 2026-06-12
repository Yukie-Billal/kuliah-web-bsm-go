package models

import (
	"github.com/goravel/framework/database/orm"
)

type AlatMusik struct {
	orm.Model
	StudioID   uint    `gorm:"column:studio_id" json:"studio_id"`
	NamaAlat   string  `gorm:"column:nama_alat" json:"nama_alat"`
	Kondisi    string  `gorm:"column:kondisi;default:Baik" json:"kondisi"`
	Keterangan *string `gorm:"column:keterangan" json:"keterangan"`

	// Relationships
	Studio     *Studio `gorm:"foreignKey:StudioID" json:"studio,omitempty"`
}

func (a *AlatMusik) TableName() string {
	return "alat_musiks"
}
