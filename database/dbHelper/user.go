package dbHelper

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"new_restaurant/database"
	"new_restaurant/models"
)

func IsUserExists(email string) (bool, error) {
	var id string
	err := database.Rest.Get(&id, `SELECT id FROM users WHERE email = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateUser(tx *sqlx.Tx, user models.User, createdBy string) (string, error) {
	var id string
	err := tx.QueryRowx(`
		INSERT INTO users (name, email, password, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		user.Name, user.Email, user.Password, createdBy,
	).Scan(&id)

	return id, err
}

func CreateUserRole(tx *sqlx.Tx, role models.UserRole) error {
	_, err := tx.NamedExec(`
		INSERT INTO user_role (user_id, role_type) 
		VALUES (:user_id, :role_type)`, &role)
	return err
}

func CreateUserAddress(address models.UserAddress) error {
	_, err := database.Rest.NamedExec(`
		INSERT INTO user_address ( user_id, address, is_primary, latitude, longitude) 
		VALUES ( :user_id, :address, :is_primary, :latitude, :longitude)`, &address)
	return err
}

//func GetUserByEmail(db *sqlx.DB, email string) (models.User, error) {
//	var user models.User
//	err := db.Get(&user, "SELECT id, name, email, password FROM users WHERE email = $1 AND archived_at IS NULL", email)
//	return user, err
//}
//
//func GetUserRoleByUserID(db *sqlx.DB, userID string) (models.UserRole, error) {
//	var role models.UserRole
//	err := database.Rest.Get(&role, "SELECT role_type FROM user_role WHERE user_id = $1 AND archived_at IS NULL", userID)
//	return role, err
//}

func GetUserByEmailAndRole(db *sqlx.DB, email string, role string) (models.User, error) {
	var user models.User
	err := db.Get(&user, `SELECT u.id, u.name, u.email,u.password, ur.role_type
		FROM users u
		JOIN user_role ur ON u.id = ur.user_id
		WHERE u.email = $1 AND ur.role_type = $2`, email, role)
	return user, err
}

func ListAllSubAdmins(db *sqlx.DB) ([]models.UserResponse, error) {
	const query = `
		SELECT u.id, u.name, u.email, ur.role_type
		FROM users u
		JOIN user_role ur ON u.id = ur.user_id
		WHERE ur.role_type = 'sub_admin' AND u.archived_at IS NULL
		LIMIT 5 OFFSET 0;`

	var subAdmins []models.UserResponse
	err := db.Select(&subAdmins, query)
	return subAdmins, err
}

func ListAllUsers(db *sqlx.DB) ([]models.UserResponse, error) {
	const query = `
		SELECT u.id, u.name, u.email, ur.role_type
		FROM users u
		JOIN user_role ur ON u.id = ur.user_id
		WHERE u.archived_at IS NULL
		LIMIT 5 OFFSET 0;`

	var user []models.UserResponse
	err := db.Select(&user, query)
	return user, err
}
