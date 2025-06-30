package dbHelper

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"new_restaurant/models"
)

func IsUserExists(db *sqlx.DB, email string) (bool, error) {
	var id string
	err := db.Get(&id, `SELECT id FROM users WHERE email = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateUser(tx *sqlx.Tx, user models.User) (string, error) {
	var id string
	err := tx.QueryRowx(`
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id`,
		user.Name, user.Email, user.Password,
	).Scan(&id)

	return id, err
}

func CreateUserRole(tx *sqlx.Tx, role models.UserRole) error {
	_, err := tx.NamedExec(`
		INSERT INTO user_role (user_id, role_type) 
		VALUES (:user_id, :role_type)`, &role)
	return err
}

func CreateUserAddress(db *sqlx.DB, address models.UserAddress) error {
	_, err := db.NamedExec(`
		INSERT INTO user_address ( user_id, address, latitude, longitude) 
		VALUES ( :user_id, :address, :latitude, :longitude)`, &address)
	return err
}

func GetUserByEmail(db *sqlx.DB, email string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT id, name, email, password FROM users WHERE email = $1 AND archived_at IS NULL", email)
	return user, err
}

func GetUserRoleByUserID(db *sqlx.DB, userID string) (models.UserRole, error) {
	var role models.UserRole
	err := db.Get(&role, "SELECT role_type FROM user_role WHERE user_id = $1 AND archived_at IS NULL", userID)
	return role, err
}

func CreateSession(db *sqlx.DB, session models.Session) error {
	_, err := db.NamedExec(`INSERT INTO user_session ( user_id, refresh_token)
        VALUES (:user_id, :refresh_token)`, &session)
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
