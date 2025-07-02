package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	if req.Name == " " || req.Address == "" || req.Rating > 5 || req.Rating < 0 {
		http.Error(w, "invalid input values", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	restaurant := models.Restaurant{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Rating:    req.Rating,
		CreatedBy: claims.UserID,
	}

	if err := dbHelper.CreateRestaurant(restaurant); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "error creating restaurant", http.StatusInternalServerError)
		return
	}

	utils.JSON.NewEncoder(w).Encode(map[string]string{"message": "restaurant created successfully"})
}

func ListAllRestaurantBySubAdmin(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") && !utils.HasRole(r, "sub_admin") {
		http.Error(w, "only admin or sub admins can access", http.StatusForbidden)
		return
	}

	restaurants, err := dbHelper.ListAllRestaurantBySubAdmin()
	if err != nil {
		http.Error(w, "failed to list restaurant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"restaurants": restaurants,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func ListAllRestaurantByAdmin(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") {
		http.Error(w, "only admin can access", http.StatusForbidden)
		return
	}

	restaurants, err := dbHelper.ListAllRestaurant()
	if err != nil {
		http.Error(w, "failed to list restaurant", http.StatusInternalServerError)
		return
	}

	// JSON Response
	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"restaurants": restaurants,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func ListAllRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurants, err := dbHelper.ListAllRestaurant()
	if err != nil {
		http.Error(w, "failed to list restaurant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"restaurants": restaurants,
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

	if req.RestaurantID == "" || req.Name == "" {
		http.Error(w, "missing inputs", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Build dish object
	dish := models.Dish{
		RestaurantID: req.RestaurantID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CreatedBy:    claims.UserID,
	}

	err := dbHelper.CreateDish(dish)
	if err != nil {
		log.Println("error creating dish:", err)
		http.Error(w, "error creating dish: ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "Dish created successfully",
	})
}

func ListAllDishByRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID := r.URL.Query().Get("id")
	if restaurantID == "" {
		http.Error(w, "restaurant ID is required", http.StatusBadRequest)
		return
	}

	dishes, err := dbHelper.ListAllDishByRestaurant(restaurantID)
	if len(dishes) == 0 {
		http.Error(w, "no dish found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("error listing dish:", err)
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
	restaurantID := r.URL.Query().Get("id")

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user address
	userAddress, err := dbHelper.GetPrimaryAddressByUserID(claims.UserID)
	if err != nil {
		log.Println("error getting user:", err)
		http.Error(w, "User address not found", http.StatusNotFound)
		return
	}

	// Get restaurant
	restaurant, err := dbHelper.GetRestaurantByID(restaurantID)
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
