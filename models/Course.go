package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name        string
	Capacity    int
	Description string
	Teacher     User
	Students    []User `gorm:"many2many:student_courses;joinForeignKey:CourseID;JoinReferences:StudentID"`
	Assignments []Assignment
}
