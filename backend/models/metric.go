package models

import "time"

type Metric struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ServerID    uint      `gorm:"index;not null" json:"server_id"`
	CPUPercent  float64   `json:"cpu_percent"`
	RAMTotal    uint64    `json:"ram_total"`
	RAMUsed     uint64    `json:"ram_used"`
	RAMPercent  float64   `json:"ram_percent"`
	DiskTotal   uint64    `json:"disk_total"`
	DiskUsed    uint64    `json:"disk_used"`
	DiskPercent float64   `json:"disk_percent"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}
