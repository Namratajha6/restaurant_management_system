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

	// Auth routes
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/GetDishesByID", handlers.ListAllDishByRestaurant).Methods("GET")
	r.HandleFunc("/GetRestaurants", handlers.ListAllRestaurant).Methods("GET")

	// Protected routes (with auth middleware)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/CreateAddress", handlers.CreateAddress).Methods("POST")
	protected.HandleFunc("/CalculateDistance", handlers.CalculateDistance).Methods("POST")

	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/CreateUser", handlers.CreateUser).Methods("POST")
	admin.HandleFunc("/GetUsers", handlers.ListAllUsers).Methods("GET")
	admin.HandleFunc("/GetSubadmins", handlers.ListAllSubAdmins).Methods("GET")
	admin.HandleFunc("/CreateRestaurants", handlers.CreateRestaurant).Methods("POST")
	admin.HandleFunc("/GetRestaurants", handlers.ListAllRestaurantByAdmin).Methods("GET")

	admin.HandleFunc("/CreateDish", handlers.CreateDish).Methods("POST")

	subAdmin := protected.PathPrefix("/subAdmin").Subrouter()
	subAdmin.HandleFunc("/GetRestaurants", handlers.ListAllRestaurantBySubAdmin).Methods("GET")

	return r
}
