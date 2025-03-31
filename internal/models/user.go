package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// // User model
// type User struct {
// 	UserID    uuid.UUID `json:"user_id" db:"user_id" redis:"user_id" validate:"omitempty"`
// 	Username  string    `json:"username" db:"username" redis:"username" validate:"required,lte=30"`
// 	Password  string    `json:"password,omitempty" db:"password" redis:"password" validate:"omitempty,required,gte=6"`
// 	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
// }

// User model
type User struct {
	UserID    uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;primaryKey;default:gen_random_uuid()" db:"user_id" redis:"user_id" validate:"omitempty"`
	Username  string    `json:"username" gorm:"column:username;type:varchar(30);not null" db:"username" redis:"username" validate:"required,lte=30"`
	Password  string    `json:"password,omitempty" gorm:"column:password;type:varchar(255);not null" db:"password" redis:"password" validate:"omitempty,required,gte=6"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at;default:current_timestamp" db:"created_at" redis:"created_at"`
}

// Find user
type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Username = strings.ToLower(strings.TrimSpace(u.Username))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}
