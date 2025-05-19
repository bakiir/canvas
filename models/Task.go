package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Title       string
	Description string
	Deadline    time.Time

	CourseID uint   `json:"course_id"`
	Course   Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
