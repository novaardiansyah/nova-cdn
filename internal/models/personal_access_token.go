package models

import "time"

type PersonalAccessToken struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	TokenableType string     `json:"tokenable_type"`
	TokenableID   uint       `json:"tokenable_id"`
	Name          string     `json:"name"`
	Token         string     `json:"token"`
	Abilities     string     `json:"abilities"`
	LastUsedAt    *time.Time `json:"last_used_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}
