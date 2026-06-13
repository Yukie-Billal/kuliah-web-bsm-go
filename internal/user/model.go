package user

import "time"

type User struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleId          string     `json:"role_id"`
	Username        string     `json:"username"`
	Email           string     `gorm:"unique;column:email;type:varchar;size:255" json:"email"`
	Password        string     `json:"-"`
	GoogleId        *string    `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	Avatar          *string    `json:"avatar"`
	RememberToken   *string    `json:"-"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
}
