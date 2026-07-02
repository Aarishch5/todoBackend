package models

type Error struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
