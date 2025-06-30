package models

import (
	"time"
)

type Restaurant struct {
	ID         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Address    string     `json:"address" db:"address"`
	Latitude   *float64   `json:"latitude,omitempty" db:"latitude"`
	Longitude  *float64   `json:"longitude,omitempty" db:"longitude"`
	CreatedBy  string     `json:"created_by" db:"created_by"`
	Rating     float64    `json:"rating" db:"rating"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type Dish struct {
	ID           string     `json:"id" db:"id"`
	RestaurantID string     `json:"restaurant_id" db:"restaurant_id"`
	Name         string     `json:"name" db:"name"`
	Description  *string    `json:"description,omitempty" db:"description"`
	Price        *float64   `json:"price,omitempty" db:"price"`
	CreatedBy    string     `json:"created_by" db:"created_by"`
	CreatedAt    *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt   *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

// RestaurantWithDishes combines Restaurant with its dishes
type RestaurantWithDishes struct {
	Restaurant Restaurant `json:"restaurant"`
	Dishes     []Dish     `json:"dishes"`
}

// CreateRestaurantRequest for API requests
type CreateRestaurantRequest struct {
	Name      string   `json:"name" validate:"required"`
	Address   string   `json:"address" validate:"required"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Rating    float64  `json:"rating" validate:"required,min=0,max=5"`
}

// UpdateRestaurantRequest for API requests
type RestaurantResponse struct {
	Name      *string  `json:"name,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Rating    *float64 `json:"rating,omitempty" validate:"omitempty,min=0,max=5"`
}

// CreateDishRequest for API requests
type CreateDishRequest struct {
	RestaurantID string   `json:"restaurant_id" validate:"required,uuid"`
	Name         string   `json:"name" validate:"required"`
	Description  *string  `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
}

// UpdateDishRequest for API requests
type UpdateDishRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
}

// RestaurantSearchRequest for search functionality
type RestaurantSearchRequest struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Radius    *float64 `json:"radius,omitempty"` // in kilometers
	MinRating *float64 `json:"min_rating,omitempty"`
	MaxRating *float64 `json:"max_rating,omitempty"`
	Limit     *int     `json:"limit,omitempty"`
	Offset    *int     `json:"offset,omitempty"`
}
