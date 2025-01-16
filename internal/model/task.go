package model

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" binding:"required,oneof=pending in-progress completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
