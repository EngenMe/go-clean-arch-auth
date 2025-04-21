package entity

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email" gorm:"uniqueIndex"`
	Password  string    `json:"-" validate:"required,min=8"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user *User) isEntity() {}
