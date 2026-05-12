package model

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus int

const (
	UserStatusActive    UserStatus = 1
	UserStatusSuspended UserStatus = -1
)

type UserRole int

const (
	RoleCustomer    UserRole = 0
	RoleAdmin       UserRole = 1
	RoleDeliveryMan UserRole = 2
	RoleSeller      UserRole = 3
)

type User struct {
	ID        uint           `gorm:"primaryKey"              json:"id"`
	Name      string         `gorm:"not null"                json:"name"`
	Surname   string         `gorm:"not null"                json:"surname"`
	Username  string         `gorm:"uniqueIndex;not null"    json:"username"`
	Email     string         `gorm:"uniqueIndex;not null"    json:"email"`
	Password  string         `gorm:"not null"                json:"-"`
	Gender    string         `                               json:"gender"`
	Dob       *time.Time     `                               json:"dob,omitempty"`
	Status    UserStatus     `gorm:"default:1"               json:"status"`
	Role      UserRole       `gorm:"default:0"               json:"role"`
	CreatedAt time.Time      `                               json:"created_at"`
	UpdatedAt time.Time      `                               json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                   json:"-"`
}
