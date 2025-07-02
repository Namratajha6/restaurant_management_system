package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"new_restaurant/handlers"
	"new_restaurant/middleware"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
			logrus.Errorf("Error encoding health response: %v", err)
		}
	}).Methods("GET")

	r.HandleFunc("/api/v1/auth/login", handlers.LoginHandler).Methods("POST")

	// Public Restaurant Info (anyone can view)
	r.HandleFunc("/api/v1/restaurants", handlers.ListAllRestaurant).Methods("GET")
	r.HandleFunc("/api/v1/restaurants/{id}/dishes", handlers.ListAllDishByRestaurant).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// User actions
	protected.HandleFunc("/users/address", handlers.CreateAddress).Methods("POST")
	protected.HandleFunc("/restaurants/{id}/distance", handlers.CalculateDistance).Methods("GET")

	// Admin Routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	admin.HandleFunc("/users", handlers.ListAllUsers).Methods("GET")
	admin.HandleFunc("/subadmins", handlers.ListAllSubAdmins).Methods("GET")
	admin.HandleFunc("/restaurants", handlers.CreateRestaurant).Methods("POST")
	admin.HandleFunc("/restaurants", handlers.ListAllRestaurantByAdmin).Methods("GET")
	admin.HandleFunc("/dishes", handlers.CreateDish).Methods("POST")

	// Subadmin Routes
	subAdmin := protected.PathPrefix("/subadmin").Subrouter()
	subAdmin.HandleFunc("/restaurants", handlers.ListAllRestaurantBySubAdmin).Methods("GET")

	return r
}
