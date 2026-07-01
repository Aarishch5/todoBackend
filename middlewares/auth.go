package middleware

import (
	"ToDo/database/dbHelper"
	"context"
	"net/http"
	"strings"

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

		sessionToken := strings.TrimPrefix(authHeader, "Bearer ")

		// session token IS the credential now — look it up directly
		session, err := dbHelper.GetSessionByToken(sessionToken)
		if err != nil {
			http.Error(w, "session revoked or expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, session.UserID)
		ctx = context.WithValue(ctx, SessionContextKey, sessionToken)

		// pass the request forward
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
