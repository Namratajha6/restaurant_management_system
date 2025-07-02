package models

import (
	"time"
)

type Restaurant struct {
	ID         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Address    string     `json:"address" db:"address"`
	Latitude   *float64   `json:"latitude" db:"latitude"`
	Longitude  *float64   `json:"longitude" db:"longitude"`
	CreatedBy  string     `json:"created_by" db:"created_by"`
	Rating     float64    `json:"rating" db:"rating"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type Dish struct {
	ID           string     `json:"id" db:"id"`
	RestaurantID string     `json:"restaurant_id" db:"restaurant_id"`
	Name         string     `json:"name" db:"name"`
	Description  *string    `json:"description" db:"description"`
	Price        *float64   `json:"price,omitempty" db:"price"`
	CreatedBy    string     `json:"created_by" db:"created_by"`
	CreatedAt    *time.Time `json:"created_at" db:"created_at"`
	ArchivedAt   *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type CreateRestaurantRequest struct {
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Rating    float64  `json:"rating"`
}

type CreateDishRequest struct {
	RestaurantID string   `json:"restaurant_id" `
	Name         string   `json:"name" `
	Description  *string  `json:"description"`
	Price        *float64 `json:"price"`
}
