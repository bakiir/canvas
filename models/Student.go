package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Name     string                 `json:"name,omitempty"`
	Login    string                 `json:"login,omitempty" gorm:"unique"`
	Password string                 `json:"-"`
	Courses  []Course               `gorm:"many2many:student_courses;" json:"courses,omitempty"`
	Grades   map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"grades,omitempty"`
}
