package models

import (
	"time"
)

type RoleType string

const (
	RoleAdmin    RoleType = "admin"
	RoleSubAdmin RoleType = "sub_admin"
	RoleUser     RoleType = "user"
)

type User struct {
	ID         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Email      string     `json:"email" db:"email"`
	RoleType   *RoleType  `json:"roleType" db:"role_type"`
	Password   string     `json:"-" db:"password"`
	CreatedBy  string     `json:"createdBy" db:"created_by"`
	CreatedAt  *time.Time `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

type UserRole struct {
	ID         string     `json:"id" db:"id"`
	UserID     string     `json:"userId" db:"user_id"`
	RoleType   RoleType   `json:"roleType" db:"role_type"`
	CreatedAt  *time.Time `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

type UserAddress struct {
	ID         string     `json:"id" db:"id"`
	UserID     string     `json:"userId" db:"user_id"`
	Address    string     `json:"address" db:"address"`
	IsPrimary  bool       `json:"isPrimary" db:"is_primary"`
	Latitude   *float64   `json:"latitude,omitempty" db:"latitude"`
	Longitude  *float64   `json:"longitude,omitempty" db:"longitude"`
	CreatedAt  *time.Time `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

// CreateUserRequest for API requests
type UserRequest struct {
	Name     string     `json:"name"`
	Email    string     `json:"email" `
	Password string     `json:"password"`
	Roles    []RoleType `json:"roles"`
}

// LoginRequest for authentication
type LoginRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	RoleType RoleType `json:"roleType"`
}

type UserResponse struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Email     string `db:"email" json:"email"`
	RoleTypes string `db:"role_type" json:"role_type"`
}

type UserAddressRequest struct {
	Name      string   `json:"name" `
	Address   string   `json:"address" `
	IsPrimary bool     `json:"isPrimary" `
	Latitude  *float64 `json:"latitude" `
	Longitude *float64 `json:"longitude" `
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
