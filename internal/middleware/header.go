package middleware

import (
	"net/http"
	"time"
)

// DefaultHeader for set default header
func DefaultHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, PUT, POST, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Accept-ranges", "items")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	now := time.Now()
	w.Header().Set("Date", now.String())

	next(w, r)
}
