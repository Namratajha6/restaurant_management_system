package dbHelper

import (
	"new_restaurant/database"
	"new_restaurant/models"
)

func CreateRestaurant(restaurant models.Restaurant) error {
	query := `INSERT INTO restaurant (name, address, latitude, longitude, created_by, rating) 
				VALUES ( :name, :address, :latitude, :longitude, :created_by, :rating)`
	_, err := database.Rest.NamedExec(query, restaurant)
	return err
}

func CreateDish(dish models.Dish) error {
	query := `INSERT INTO dishes ( restaurant_id, name, description, price, created_by) 
				VALUES( :restaurant_id, :name, :description, :price, :created_by)`
	_, err := database.Rest.NamedExec(query, dish)
	return err
}

func ListAllDishByRestaurant(page int, limit int, restaurantID string) ([]models.Dish, error) {
	const query = `
		SELECT id, restaurant_id, name, description, price, created_by
		FROM dishes
		WHERE restaurant_id = $1 AND archived_at IS NULL
		LIMIT $2 OFFSET $3;`

	dishes := make([]models.Dish, 0)
	err := database.Rest.Select(&dishes, query, restaurantID, limit, limit*(page-1))
	return dishes, err
}

func ListAllRestaurant(page int, limit int) ([]models.Restaurant, error) {
	const query = `
		SELECT ID,name, address, latitude, longitude, created_by, rating
		FROM restaurant
		WHERE archived_at IS NULL
		LIMIT $1 OFFSET $2 ;`

	offset := limit * (page - 1)
	restaurants := make([]models.Restaurant, 0)
	err := database.Rest.Select(&restaurants, query, limit, offset)
	return restaurants, err
}

func ListAllRestaurantBySubAdmin(page int, limit int) ([]models.Restaurant, error) {
	const query = `
		SELECT r.id, r.name, r.address, r.latitude, r.longitude, r.created_by, r.rating
		FROM restaurant r
		JOIN user_role ur ON r.created_by = ur.user_id
		WHERE ur.role_type = 'sub_admin' AND r.archived_at IS NULL
		LIMIT $1 OFFSET $2;`

	restaurants := make([]models.Restaurant, 0)
	err := database.Rest.Select(&restaurants, query, limit, limit*(page-1))
	return restaurants, err
}

func GetRestaurantByID(restaurantID string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	query := `SELECT id, latitude, longitude
	          FROM restaurant 
	          WHERE id = $1 AND archived_at IS NULL`
	err := database.Rest.Get(&restaurant, query, restaurantID)
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func GetPrimaryAddressByUserID(userID string) (*models.UserAddress, error) {
	var address models.UserAddress
	query := `
		SELECT id, user_id, address, latitude, longitude
		FROM user_address
		WHERE user_id = $1 AND is_primary = TRUE AND archived_at IS NULL;
	`
	err := database.Rest.Get(&address, query, userID)
	if err != nil {
		return nil, err
	}
	return &address, nil
}
