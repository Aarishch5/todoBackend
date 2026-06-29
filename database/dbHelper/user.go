package dbHelper

import (
	"ToDo/database/migrations"
	"ToDo/models"

	"github.com/google/uuid"
)

func CreateUser(user *models.User) error {

	query := `
	INSERT INTO users(name,email,password)
	VALUES($1,$2,$3)
	RETURNING user_id
	`

	return migrations.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.UserID)
}

func DeleteUser(userID uuid.UUID) (int64, error) {

	result, err := migrations.DB.Exec(
		`DELETE FROM users WHERE user_id=$1`,
		userID,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func GetUserByEmail(email string) (models.User, error) {

	var user models.User

	err := migrations.DB.Get(
		&user,
		`
		SELECT user_id,name,email,password
		FROM users
		WHERE email=$1
		`,
		email,
	)

	return user, err
}
