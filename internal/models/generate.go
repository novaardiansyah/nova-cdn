package models

import (
	"time"

	"gorm.io/gorm"
)

type Generate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Alias     string         `json:"alias"`
	Prefix    *string        `json:"prefix"`
	Suffix    *string        `json:"suffix"`
	Queue     int            `json:"queue"`
	Separator string         `json:"separator"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" swaggertype:"string"`
}

func (Generate) TableName() string {
	return "generates"
}
