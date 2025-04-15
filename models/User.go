package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"login,omitempty" gorm:"unique"`
	Password string   `json:"-"`
	Role     string   `json:"role,omitempty"`
	Courses  []Course `gorm:"many2many:student_courses;joinForeignKey:StudentID;JoinReferences:CourseID"`
}
