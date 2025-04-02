package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name      string
	Capacity  int
	TeacherID uint
	Teacher   Teacher   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Students  []Student `gorm:"many2many:student_courses;"`
}
