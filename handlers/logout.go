package handlers

import (
	"ToDo/database/dbHelper"
	middleware "ToDo/middlewares"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken, ok := middleware.GetSessionToken(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if _, err := dbHelper.DeleteSession(sessionToken, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
