// models/user.go
package models

import (
	"github.com/google/uuid"
	"time"
)

type RoleType string

const (
	RoleAdmin    RoleType = "admin"
	RoleSubAdmin RoleType = "sub_admin"
	RoleUser     RoleType = "user"
)

type User struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Email      string     `json:"email" db:"email"`
	Password   string     `json:"-" db:"password"` // "-" to exclude from JSON
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type UserRole struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	RoleType   RoleType   `json:"role_type" db:"role_type"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type UserAddress struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	Address    string     `json:"address" db:"address"`
	Latitude   *float64   `json:"latitude,omitempty" db:"latitude"`
	Longitude  *float64   `json:"longitude,omitempty" db:"longitude"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

// CreateUserRequest for API requests
type UserRequest struct {
	Name     string     `json:"name" validate:"required"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=6"`
	Roles    []RoleType `json:"roles" validate:"required,min=1,dive,required"`
}

// LoginRequest for authentication
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Session struct {
	ID           uuid.UUID  `db:"id"`
	UserID       uuid.UUID  `db:"user_id"`
	RefreshToken string     `db:"refresh_token"`
	CreatedAt    *time.Time `db:"created_at"`
	ArchivedAt   *time.Time `db:"archived_at"`
}

type UserResponse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	RoleTypes string    `db:"role_type" json:"role_type"` // array
}

type UserAddressRequest struct {
	Name      string   `json:"name" validate:"required"`
	Address   string   `json:"address" validate:"required,email"`
	Latitude  *float64 `json:"latitude" validate:"required,min=0"`
	Longitude *float64 `json:"longitude" validate:"required,min=0"`
}

type DistanceRequest struct {
	UserAddressID string `json:"user_address_id"`
	RestaurantID  string `json:"restaurant_id"`
}

// represents the distance calculation response
type DistanceResponse struct {
	Distance float64 `json:"distance_km"`
	Message  string  `json:"message"`
}
