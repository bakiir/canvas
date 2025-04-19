package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Name     string
	Login    string `gorm:"unique"`
	Password string
	Courses  []Course `gorm:"foreignKey:TeacherID;"`
}

// Хэширование пароля
func (t *Teacher) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.Password = string(hashed)
	return nil
}

// Проверка пароля
func (t *Teacher) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Password), []byte(password))
	return err == nil
}
