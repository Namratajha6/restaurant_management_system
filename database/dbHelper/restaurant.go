package dbHelper

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"new_restaurant/models"
)

func CreateRestaurant(db *sqlx.DB, restaurant models.Restaurant) error {
	query := `INSERT INTO restaurant (id, name, address, latitude, longitude, created_by, rating) 
				VALUES ( :id, :name, :address, :latitude, :longitude, :created_by, :rating)`
	_, err := db.NamedExec(query, restaurant)
	return err
}

func CreateDish(db *sqlx.DB, dish models.Dish) error {
	query := `INSERT INTO dishes (id, restaurant_id, name, description, price, created_by) 
				VALUES(:id, :restaurant_id, :name, :description, :price, :created_by)`
	_, err := db.NamedExec(query, dish)
	return err
}

func ListAllDishByRestaurant(db *sqlx.DB, restaurantID uuid.UUID) ([]models.Dish, error) {
	const query = `
		SELECT id, restaurant_id, name, description, price, created_by
		FROM dishes
		WHERE restaurant_id = $1 AND archived_at IS NULL;`

	var dishes []models.Dish
	err := db.Select(&dishes, query, restaurantID)
	return dishes, err
}

func ListAllRestaurant(db *sqlx.DB) ([]models.Restaurant, error) {
	const query = `
		SELECT ID,name, address, latitude, longitude, created_by, rating
		FROM restaurant
		WHERE archived_at IS NULL;`

	var restaurant []models.Restaurant
	err := db.Select(&restaurant, query)
	return restaurant, err
}

func ListAllRestaurantBySubAdmin(db *sqlx.DB) ([]models.Restaurant, error) {
	const query = `
		SELECT r.id, r.name, r.address, r.latitude, r.longitude, r.created_by, r.rating
		FROM restaurant r
		JOIN user_role ur ON r.created_by = ur.user_id
		WHERE ur.role_type = 'sub_admin' AND r.archived_at IS NULL;`

	var restaurants []models.Restaurant
	err := db.Select(&restaurants, query)
	return restaurants, err
}

func GetRestaurantByID(db *sqlx.DB, restaurantID string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	query := `SELECT id, name, address, latitude, longitude, rating, created_by
	          FROM restaurant 
	          WHERE id = $1 AND archived_at IS NULL`
	err := db.Get(&restaurant, query, restaurantID)
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func GetUserAddress(db *sqlx.DB, addressID string) (*models.UserAddress, error) {
	var address models.UserAddress
	query := `SELECT id, user_id, address, latitude, longitude
	          FROM user_address 
	          WHERE id = $1 `
	err := db.Get(&address, query, addressID)
	if err != nil {
		return nil, err
	}
	return &address, nil
}
