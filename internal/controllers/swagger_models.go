package controllers

import "time"

type GallerySwagger struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	Url          string    `json:"url"`
	FileSize     uint32    `json:"file_size"`
	IsPrivate    bool      `json:"is_private"`
	Description  string    `json:"description"`
	Size         string    `json:"size"`
	HasOptimized bool      `json:"has_optimized"`
	GroupCode    string    `json:"group_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *string   `json:"deleted_at"`
}
