package models

import (
	"github.com/google/uuid"
	"time"
	// "golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uuid.UUID    `json:"id" gorm:"column:id; primary_key" swaggerignore:"true"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	Email       string       `json:"email" gorm:"unique" binding:"required" validate:"required,email"`
	Password    string       `json:"-" gorm:"column:password" binding:"required" validate:"required,min=6,max=20"`
	CreatedAt   time.Time    `json:"created_at" default:"current_timestamp"`
	UpdatedAt   time.Time    `json:"updated_at" default:"current_timestamp"`
} //@name User


type Login struct {
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
} //@name Login

type Register struct {
	FirstName string `json:"first_name" binding:"required" validate:"required,min=2,max=20"`
	LastName  string `json:"last_name" binding:"required" validate:"required,min=2,max=20"`
	Email     string `json:"email" binding:"required" validate:"required,email"`
	Password  string `json:"password" binding:"required" validate:"required,min=6,max=20"`
} //@name Register

