package model

import "time"

type PostgresUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Account   string    `gorm:"uniqueIndex" json:"account"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfiles struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Account   string    `gorm:"uniqueIndex" json:"account"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PostgresUser) TableName() string {
	return "users"
}

func (UserProfiles) TableName() string {
	return "user_profiles"
} 