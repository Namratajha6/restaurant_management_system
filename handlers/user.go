package handlers

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	userID := uuid.New()
	user := models.User{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		if err := dbHelper.CreateUser(tx, user); err != nil {
			return err
		}

		for _, role := range req.Roles {
			roleEntry := models.UserRole{
				ID:       uuid.New(),
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

	w.Header().Set("Content-Type", "application/json")
	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"user_id": userID.String(),
	})
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	// Only admins allowed
	if !utils.HasRole(r, "admin") {
		http.Error(w, "only admin can view sub-admins", http.StatusForbidden)
		return
	}

	// Fetch from DB
	users, err := dbHelper.ListAllUsers(database.Rest)
	if err != nil {
		http.Error(w, "failed to list sub-admins", http.StatusInternalServerError)
		return
	}

	// JSON Response
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

	user, err := dbHelper.GetUserByEmail(database.Rest, req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	role, err := dbHelper.GetUserRoleByUserID(database.Rest, user.ID)
	if err != nil {
		http.Error(w, "failed to fetch user role", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), string(role.RoleType))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), string(role.RoleType))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	session := models.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}

	err = dbHelper.CreateSession(database.Rest, session)
	if err != nil {
		http.Error(w, "failed to create refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	var req reqBody
	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := dbHelper.DeleteSessionByToken(database.Rest, req.RefreshToken)
	if err != nil {
		http.Error(w, "failed to delete session", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	utils.JSON.NewEncoder(w).Encode(map[string]string{"message": "logged out successfully"})
}

func ListAllSubAdmins(w http.ResponseWriter, r *http.Request) {
	// Only admins allowed
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

	// JSON Response
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

	userID, ok := utils.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := models.UserAddress{
		ID:        uuid.New(),
		UserID:    userID,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	err := dbHelper.CreateUserAddress(database.Rest, user)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "User address created successfully",
	})
}
