package models

import (
	"gorm.io/gorm"
	"time"
)

type Homework struct {
	gorm.Model
	StudentID uint
	TaskID    uint

	FileURL    string // сюда кладёшь путь/ссылку от файл-сервиса
	UploadedAt time.Time

	Student Student `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Task    Task    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
