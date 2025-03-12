package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgressConnection(cs string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cs), &gorm.Config{})
}
