package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name     string                 `json:"name,omitempty"`
	Login    string                 `json:"login,omitempty" gorm:"unique"`
	Password string                 `json:"-"`
	Courses  []Course               `gorm:"many2many:student_courses;" json:"courses,omitempty"`
	Grades   map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"grades,omitempty"`
}

func (t *Student) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.Password = string(hashed)
	return nil
}

func (t *Student) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Password), []byte(password))
	return err == nil
}
