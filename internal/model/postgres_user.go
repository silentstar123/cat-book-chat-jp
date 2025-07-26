package model

import "time"

type PostgresUser struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Account    string     `gorm:"uniqueIndex" json:"account"`
	Nickname   string     `json:"nickname"`
	Role       string     `json:"role"`
	Phone      *string    `json:"phone"`
	Email      *string    `json:"email"`
	Password   string     `json:"password"`
	Status     int64      `gorm:"not null;default:1" json:"status"`
	EmailValid bool       `gorm:"column:emailValid;not null;default:false" json:"emailValid"`
	PhoneValid bool       `gorm:"column:phoneValid;not null;default:false" json:"phoneValid"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
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
