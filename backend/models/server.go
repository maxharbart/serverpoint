package models

import "time"

type Server struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:200;not null" json:"name"`
	Host              string    `gorm:"size:200;not null" json:"host"`
	Port              int       `gorm:"default:22;not null" json:"port"`
	Username          string    `gorm:"size:100;not null" json:"username"`
	AuthType          string    `gorm:"size:20;not null" json:"auth_type"` // "password" or "key"
	EncryptedPassword string    `gorm:"type:text" json:"-"`
	EncryptedKey      string    `gorm:"type:text" json:"-"`
	OS                string    `gorm:"size:200" json:"os"`
	Hostname          string    `gorm:"size:200" json:"hostname"`
	Uptime            string    `gorm:"size:200" json:"uptime"`
	Online            bool      `gorm:"default:false" json:"online"`
	UserID            uint      `gorm:"not null" json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
