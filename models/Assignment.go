package models

import (
	"time"
)

type Assignment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`

	CourseID uint   `json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"-"`
}
