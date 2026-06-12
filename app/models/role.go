package models

import (
	"github.com/goravel/framework/database/orm"
)

type Role struct {
	orm.Model
	Role  string `gorm:"column:role" json:"role"`
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

func (r *Role) TableName() string {
	return "roles"
}
