package controllers

import "time"

type GallerySwagger struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	FileName    string    `json:"file_name"`
	Url         string    `json:"url"`
	FileSize    uint32    `json:"file_size"`
	IsPrivate   bool      `json:"is_private"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *string   `json:"deleted_at"`
}
