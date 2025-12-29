package models

import (
	"nova-cdn/internal/config"
	"time"

	"gorm.io/gorm"
)

type Gallery struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `json:"user_id"`
	FileName    string         `json:"file_name"`
	FilePath    string         `json:"-"`
	Url         string         `gorm:"-" json:"url"`
	FileSize    uint32         `json:"file_size"`
	IsPrivate   bool           `json:"is_private"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" swaggertype:"string"`
}

func (Gallery) TableName() string {
	return "galleries"
}

func (g *Gallery) AfterFind(tx *gorm.DB) error {
	if g.FilePath != "" {
		g.Url = config.AppURL + "/" + g.FilePath
	}
	return nil
}
