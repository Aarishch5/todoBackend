package dbHelper

import (
	"ToDo/database/migrations"
	"ToDo/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func CreateSession(session *models.Session) error {
	query := `
	INSERT INTO user_session(user_id, session_token, expires_at)
	VALUES($1,$2, $3)
	RETURNING id, created_at
	`

	return migrations.DB.QueryRow(
		query, session.UserID, session.SessionToken, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func GetSessionByToken(token string) (models.Session, error) {
	var session models.Session

	err := migrations.DB.Get(
		&session,
		`SELECT id, user_id, session_token, created_at, expires_at
		 FROM user_session
		 WHERE session_token=$1 AND expires_at > now() AND archived_at IS NULL`,
		token,
	)

	return session, err
}

func DeleteSession(token string, userID uuid.UUID) (int64, error) {
	result, err := migrations.DB.Exec(
		`UPDATE user_session SET archived_at=now() WHERE session_token=$1 AND user_id=$2`,
		token, userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func CreateSessionTx(tx *sqlx.Tx, session *models.Session) error {
	query := `
	INSERT INTO user_session(user_id, session_token, expires_at)
	VALUES($1,$2,$3)
	RETURNING id, created_at
	`
	return tx.QueryRow(query, session.UserID, session.SessionToken, session.ExpiresAt).
		Scan(&session.ID, &session.CreatedAt)
}

func DeleteSessionById(userID uuid.UUID) (int64, error) {
	result, err := migrations.DB.Exec(
		`UPDATE user_session SET archived_at=now() WHERE user_id=$1`, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
