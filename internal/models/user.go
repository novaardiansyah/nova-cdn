package models

import "time"

type User struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	Name                 string     `json:"name"`
	Email                string     `json:"email" gorm:"unique"`
	Password             string     `json:"-"`
	HasAllowNotification *bool      `json:"has_allow_notification"`
	NotificationToken    *string    `json:"notification_token,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
