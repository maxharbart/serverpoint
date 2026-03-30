package models

import "time"

type Process struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ServerID  uint      `gorm:"index;not null" json:"server_id"`
	PID       int       `json:"pid"`
	User      string    `gorm:"size:100" json:"user"`
	CPU       float64   `json:"cpu"`
	RAM       float64   `json:"ram"`
	Command   string    `gorm:"size:500" json:"command"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
