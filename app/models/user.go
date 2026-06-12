package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	RoleID          uint       `gorm:"column:role_id;default:2" json:"role_id"`
	Username        string     `gorm:"column:username" json:"username"`
	Email           string     `gorm:"column:email;unique" json:"email"`
	Password        string     `gorm:"column:password" json:"-"`
	GoogleID        *string    `gorm:"column:google_id;unique" json:"google_id"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at" json:"email_verified_at"`
	Avatar          *string    `gorm:"column:avatar" json:"avatar"`
	RememberToken   *string    `gorm:"column:remember_token" json:"-"`

	// Relationships
	Role     *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Customer *Customer `gorm:"foreignKey:UserID" json:"customer,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}
