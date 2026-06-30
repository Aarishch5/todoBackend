package handlers

import (
	"ToDo/database/dbHelper"
	"ToDo/models"
	"ToDo/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validating the name, if the name is an empty string
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// validating the email
	if err := utils.ValidateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validating password
	if err := utils.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
	}

	if err := dbHelper.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// creating a new session row in db
	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session := models.Session{
		UserID:       user.UserID,
		SessionToken: sessionToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	if err := dbHelper.CreateSession(&session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//passing the sessionToken into the jwt
	token, err := utils.GenerateJWTToken(user.UserID, sessionToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"user":  user,
		"token": token,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := dbHelper.GetUserByEmail(req.Email)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := utils.CheckPassword(
		user.Password,
		req.Password,
	); err != nil {

		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// creating a session row
	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session := models.Session{
		UserID:       user.UserID,
		SessionToken: sessionToken,
	}
	if err := dbHelper.CreateSession(&session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// passing sessionToken into the JWT
	token, err := utils.GenerateJWTToken(user.UserID, sessionToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
