package utils

import (
	"encoding/json"
	"net/http"

	"ToDo/logger"
	"ToDo/models"
)

func EncodeBody(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func RespondError(w http.ResponseWriter, statusCode int, err error, message string) {
	if err != nil {
		logger.Log.WithError(err).WithField("status_code", statusCode).Error(message)
	}

	errStr := ""
	if err != nil && statusCode < http.StatusInternalServerError {
		errStr = err.Error()
	}

	newErr := models.Error{Error: errStr, StatusCode: statusCode, Message: message}
	if encErr := EncodeBody(w, statusCode, newErr); encErr != nil {
		logger.Log.WithError(encErr).Error("failed to encode error response")
	}
}
