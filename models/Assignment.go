package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	Title       string
	Description string
	Courses     Course
	Deadline    pgtype.Date
}
