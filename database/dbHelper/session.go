package dbHelper

import (
	"ToDo/database/migrations"
	"ToDo/models"

	"github.com/google/uuid"
)

func CreateSession(session *models.Session) error {
	query := `
	INSERT INTO user_session(user_id, session_token)
	VALUES($1,$2)
	RETURNING id, created_at
	`

	return migrations.DB.QueryRow(
		query, session.UserID, session.SessionToken,
	).Scan(&session.ID, &session.CreatedAt)
}

func GetSessionByToken(token string) (models.Session, error) {
	var session models.Session

	err := migrations.DB.Get(
		&session,
		`SELECT id, user_id, session_token, created_at FROM user_session WHERE session_token=$1`,
		token,
	)

	return session, err
}

func DeleteSession(token string, userID uuid.UUID) (int64, error) {
	result, err := migrations.DB.Exec(
		`DELETE FROM user_session WHERE session_token=$1 AND user_id=$2`,
		token, userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteAllSessions(userID uuid.UUID) (int64, error) {
	result, err := migrations.DB.Exec(
		`DELETE FROM user_session WHERE user_id=$1`,
		userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
