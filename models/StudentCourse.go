package models

type StudentCourse struct {
	StudentID uint `gorm:"primaryKey"`
	CourseID  uint `gorm:"primaryKey"`
}
