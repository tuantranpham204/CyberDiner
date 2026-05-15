package entity

import "time"

type Role string

const (
	RoleUser        Role = "user"
	RoleStaff       Role = "staff"
	RoleDeliveryMan Role = "delivery_man"
	RoleAdmin       Role = "admin"
)

func (r Role) String() string { return string(r) }

var AllRoles = []Role{RoleUser, RoleStaff, RoleDeliveryMan, RoleAdmin}

type Gender string

const (
	GenderMale           Gender = "male"
	GenderFemale         Gender = "female"
	GenderOther          Gender = "other"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
)

func (g Gender) String() string { return string(g) }

var AllGenders = []Gender{GenderMale, GenderFemale, GenderOther, GenderPreferNotToSay}

type User struct {
	ID           int64     `gorm:"type:bigserial;primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(30);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         Role      `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	IsBanned     bool      `gorm:"not null;default:false" json:"is_banned"`
	Profile      *Profile  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Profile struct {
	UserID      int64     `gorm:"type:bigint;primaryKey;column:user_id" json:"user_id"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Surname     string    `gorm:"type:varchar(50);not null" json:"surname"`
	DateOfBirth time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	Gender      Gender    `gorm:"type:varchar(20);not null" json:"gender"`
	PhoneNumber *string   `gorm:"type:varchar(20)" json:"phone_number,omitempty"`
	Address     *string   `gorm:"type:text" json:"address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Profile) TableName() string { return "profiles" }
