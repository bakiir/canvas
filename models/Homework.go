package models

import (
	"gorm.io/gorm"
	"time"
)

type Homework struct {
	gorm.Model
	StudentID uint `gorm:"uniqueIndex:idx_student_task"`
	TaskID    uint `gorm:"uniqueIndex:idx_student_task"`

	FileURL    string // сюда кладёшь путь/ссылку от файл-сервиса
	UploadedAt time.Time

	Student Student `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Task    Task    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
