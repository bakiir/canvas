package models

import "time"

type StudentCourse struct {
	StudentID uint `gorm:"primaryKey"`
	CourseID  uint `gorm:"primaryKey"`
	CreatedAt time.Time
}
