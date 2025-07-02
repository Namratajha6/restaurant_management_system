package handlers

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"new_restaurant/database"
	"new_restaurant/database/dbHelper"
	"new_restaurant/models"
	"new_restaurant/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") {
		http.Error(w, "only admin can create users", http.StatusForbidden)
		return
	}

	var req models.UserRequest
	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == " " || req.Email == "" || req.Password == "" ||
		len(req.Roles) == 0 {
		http.Error(w, "missing inputs", http.StatusBadRequest)
		return
	}

	ok, err := dbHelper.IsUserExists(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ok {
		http.Error(w, "user already exist", http.StatusNotFound)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err := dbHelper.CreateUser(tx, user, claims.UserID)
		if err != nil {
			return err
		}

		for _, role := range req.Roles {
			roleEntry := models.UserRole{
				UserID:   userID,
				RoleType: role,
			}
			if err := dbHelper.CreateUserRole(tx, roleEntry); err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		http.Error(w, "failed to create user with role", http.StatusInternalServerError)
		return
	}

	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") {
		http.Error(w, "only admin can view users", http.StatusForbidden)
		return
	}

	users, err := dbHelper.ListAllUsers(database.Rest)
	if err != nil {
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" ||
		req.RoleType == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	user, err := dbHelper.GetUserByEmailAndRole(database.Rest, req.Email, string(req.RoleType))
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "invalid email or role", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, string(req.RoleType))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func ListAllSubAdmins(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") {
		http.Error(w, "only admin can view sub-admins", http.StatusForbidden)
		return
	}

	// Fetch from DB
	subadmins, err := dbHelper.ListAllSubAdmins(database.Rest)
	if err != nil {
		http.Error(w, "failed to list sub-admins", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"subadmins": subadmins,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	if !utils.HasRole(r, "admin") && !utils.HasRole(r, "sub_admin") && !utils.HasRole(r, "user") {
		http.Error(w, "user not logged in", http.StatusForbidden)
		return
	}

	var req models.UserAddressRequest
	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == " " || req.Address == "" {
		http.Error(w, "missing inputs", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := models.UserAddress{
		UserID:    claims.UserID,
		Address:   req.Address,
		IsPrimary: req.IsPrimary,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	err := dbHelper.CreateUserAddress(user)
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "failed to create user address", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "User address created successfully",
	})
}
