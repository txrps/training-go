package database

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormFromPool(pool *sql.DB) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: pool,
	}), &gorm.Config{})
}
