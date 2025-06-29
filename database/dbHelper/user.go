package dbHelper

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"new_restaurant/models"
)

func CreateUser(tx *sqlx.Tx, user models.User) error {
	_, err := tx.NamedExec(`
		INSERT INTO users (id, name, email, password) 
		VALUES (:id, :name, :email, :password)`, &user)
	return err
}

func CreateUserRole(tx *sqlx.Tx, role models.UserRole) error {
	_, err := tx.NamedExec(`
		INSERT INTO user_role (id, user_id, role_type) 
		VALUES (:id, :user_id, :role_type)`, &role)
	return err
}

func CreateUserAddress(db *sqlx.DB, address models.UserAddress) error {
	_, err := db.NamedExec(`
		INSERT INTO user_address (id, user_id, address, latitude, longitude) 
		VALUES (:id, :user_id, :address, :latitude, :longitude)`, &address)
	return err
}

func GetUserByEmail(db *sqlx.DB, email string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1 AND archived_at IS NULL", email)
	return user, err
}

func GetUserRoleByUserID(db *sqlx.DB, userID uuid.UUID) (models.UserRole, error) {
	var role models.UserRole
	err := db.Get(&role, "SELECT * FROM user_role WHERE user_id = $1 AND archived_at IS NULL", userID)
	return role, err
}

func CreateSession(db *sqlx.DB, session models.Session) error {
	_, err := db.NamedExec(`INSERT INTO user_session (id, user_id, refresh_token)
        VALUES (:id, :user_id, :refresh_token)`, &session)
	return err
}

func DeleteSessionByToken(db *sqlx.DB, refreshToken string) error {
	_, err := db.Exec(`DELETE FROM user_session WHERE refresh_token = $1`, refreshToken)
	return err
}

func GetSessionByToken(db *sqlx.DB, refreshToken string) (models.Session, error) {
	var session models.Session
	err := db.Get(&session, `SELECT * FROM user_session WHERE refresh_token = $1`, refreshToken)
	return session, err
}

func ListAllSubAdmins(db *sqlx.DB) ([]models.UserResponse, error) {
	const query = `
		SELECT u.id, u.name, u.email, ur.role_type
		FROM users u
		JOIN user_role ur ON u.id = ur.user_id
		WHERE ur.role_type = 'sub_admin' AND u.archived_at IS NULL;`

	var subAdmins []models.UserResponse
	err := db.Select(&subAdmins, query)
	return subAdmins, err
}

func ListAllUsers(db *sqlx.DB) ([]models.UserResponse, error) {
	const query = `
		SELECT u.id, u.name, u.email, ur.role_type
		FROM users u
		JOIN user_role ur ON u.id = ur.user_id
		WHERE u.archived_at IS NULL;`

	var user []models.UserResponse
	err := db.Select(&user, query)
	return user, err
}
