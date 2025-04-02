package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Name     string
	Login    string `gorm:"unique"`
	Password string
	Courses  []Course `gorm:"foreignKey:TeacherID;"`
}
