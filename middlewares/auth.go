package middleware

import (
	"ToDo/database/dbHelper"
	"context"
	"net/http"
	"strings"

	"ToDo/utils"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserContextKey    ContextKey = "user"
	SessionContextKey ContextKey = "session"
)

func GetSessionToken(r *http.Request) (string, bool) {
	token, ok := r.Context().Value(SessionContextKey).(string)
	return token, ok
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserContextKey).(uuid.UUID)
	return userID, ok
}

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// checks for the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// checks the bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)

		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusUnauthorized)
			return
		}

		// to confirm that the session still exists in the DB
		session, err := dbHelper.GetSessionByToken(claims.SessionToken)
		if err != nil {
			http.Error(w, "session revoked", http.StatusUnauthorized)
			return
		}

		if session.UserID != userID {
			http.Error(w, "session mismatch", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		ctx = context.WithValue(ctx, SessionContextKey, claims.SessionToken)

		// pass the request forward
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
