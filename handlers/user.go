package handlers

import (
	"ToDo/database/dbHelper"
	middleware "ToDo/middlewares"
	"ToDo/utils"
	"net/http"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := dbHelper.DeleteUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	rows, err = dbHelper.DeleteAllToDos(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	if rows == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	rows, err = dbHelper.DeleteSessionById(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	if rows == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
