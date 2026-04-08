package models

import "time"

type Todo struct {
	ID        int64     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title" db:"title" gorm:"type:varchar(255);not null"`
	Completed bool      `json:"completed" db:"completed" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}
