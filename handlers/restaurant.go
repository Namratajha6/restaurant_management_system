package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"new_restaurant/database"
	"new_restaurant/database/dbHelper"
	"new_restaurant/models"
	"new_restaurant/utils"
)

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") && !utils.HasRole(r, "sub_admin") {
		http.Error(w, "only admin and sub_admin can create restaurants", http.StatusForbidden)
		return
	}

	var req models.CreateRestaurantRequest

	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := utils.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	restaurant := models.Restaurant{
		ID:        uuid.New(),
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Rating:    req.Rating,
		CreatedBy: userID,
	}

	if err := dbHelper.CreateRestaurant(database.Rest, restaurant); err != nil {
		http.Error(w, "error creating restaurant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.JSON.NewEncoder(w).Encode(map[string]string{"message": "restaurant created successfully"})
}

func ListAllRestaurantBySubAdmin(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") && !utils.HasRole(r, "sub_admin") {
		http.Error(w, "only admin or sub admins can access", http.StatusForbidden)
		return
	}

	restaurant, err := dbHelper.ListAllRestaurantBySubAdmin(database.Rest)
	if err != nil {
		http.Error(w, "failed to list restaurant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"restaurants": restaurant,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func ListAllRestaurantByAdmin(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") {
		http.Error(w, "only admin can access", http.StatusForbidden)
		return
	}

	restaurant, err := dbHelper.ListAllRestaurant(database.Rest)
	if err != nil {
		http.Error(w, "failed to list restaurant", http.StatusInternalServerError)
		return
	}

	// JSON Response
	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"restaurants": restaurant,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func ListAllRestaurant(w http.ResponseWriter, r *http.Request) {

	// Fetch from DB
	restaurant, err := dbHelper.ListAllRestaurant(database.Rest)
	if err != nil {
		http.Error(w, "failed to list restaurant", http.StatusInternalServerError)
		return
	}

	// JSON Response
	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"restaurants": restaurant,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") && !utils.HasRole(r, "sub_admin") {
		http.Error(w, "only admin or subadmin can create dishes", http.StatusForbidden)
		return
	}

	var req models.CreateDishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate UUID
	restaurantUUID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		http.Error(w, "invalid restaurant_id", http.StatusBadRequest)
		return
	}

	// Parse user ID from context
	userID, ok := utils.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Build dish object
	dish := models.Dish{
		ID:           uuid.New(),
		RestaurantID: restaurantUUID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CreatedBy:    userID,
	}

	err = dbHelper.CreateDish(database.Rest, dish)
	if err != nil {
		http.Error(w, "error creating dish: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "Dish created successfully",
	})
}

func ListAllDishByRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantIDStr := r.URL.Query().Get("id")
	if restaurantIDStr == "" {
		http.Error(w, "restaurant ID is required", http.StatusBadRequest)
		return
	}

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		http.Error(w, "invalid restaurant ID format", http.StatusBadRequest)
		return
	}

	dishes, err := dbHelper.ListAllDishByRestaurant(database.Rest, restaurantID)
	if err != nil {
		http.Error(w, "failed to list dishes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"dishes": dishes,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func CalculateDistance(w http.ResponseWriter, r *http.Request) {
	var req models.DistanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user address
	userAddress, err := dbHelper.GetUserAddress(database.Rest, req.UserAddressID)
	if err != nil {
		http.Error(w, "User address not found", http.StatusNotFound)
		return
	}

	// Get restaurant
	restaurant, err := dbHelper.GetRestaurantByID(database.Rest, req.RestaurantID)
	if err != nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	// Ensure lat/long values are not nil
	if userAddress.Latitude == nil || userAddress.Longitude == nil ||
		restaurant.Latitude == nil || restaurant.Longitude == nil {
		http.Error(w, "Missing coordinates for distance calculation", http.StatusBadRequest)
		return
	}

	// Calculate distance (in KM)
	distance := utils.CalculateDistance(
		*userAddress.Latitude, *userAddress.Longitude,
		*restaurant.Latitude, *restaurant.Longitude,
	)

	// Send response
	response := models.DistanceResponse{
		Distance: distance,
		Message:  "Distance calculated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
