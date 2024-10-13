package models

import (
	"github.com/jinzhu/gorm"
)

// User represents a user model
type User struct {
	Base
	Name  string `json:"name" gorm:"not null" binding:"required"`
	Email string `json:"email" gorm:"not null;unique" binding:"required,email"`
	Age   int    `json:"age"`
}

// AdjustFieldErrors adjusts field errors to remove the model prefix and convert to lowercase
func (u *User) AdjustFieldErrors(errMap map[string][]string) map[string][]string {
	return u.Base.AdjustFieldErrors(errMap, u.ModelName())
}

// ModelName returns the name of the model
func (u *User) ModelName() string {
	return "User"
}

var DB *gorm.DB
