package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]User, error) {

	var users []User

	err := r.db.Find(&users).Error

	return users, err
}

func (r *Repository) FindByID(id int) (*User, error) {

	var user User

	err := r.db.First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
