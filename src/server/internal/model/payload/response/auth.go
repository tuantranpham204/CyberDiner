package response

import (
	"time"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/entity"
)

type User struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	IsBanned    bool       `json:"is_banned"`
	Name        string     `json:"name"`
	Surname     string     `json:"surname"`
	Gender      string     `json:"gender"`
	DOB         string     `json:"dob"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	Address     *string    `json:"address,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Auth struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func FromUser(u *entity.User) User {
	r := User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		IsBanned:  u.IsBanned,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.Profile != nil {
		r.Name = u.Profile.Name
		r.Surname = u.Profile.Surname
		r.Gender = string(u.Profile.Gender)
		r.DOB = u.Profile.DateOfBirth.Format("2006-01-02")
		r.PhoneNumber = u.Profile.PhoneNumber
		r.Address = u.Profile.Address
	}
	return r
}
