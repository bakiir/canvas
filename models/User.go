package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name,omitempty"`
	Email    string `json:"login,omitempty" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role,omitempty"`
}
