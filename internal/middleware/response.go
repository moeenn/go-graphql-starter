package middleware

import (
	"encoding/json"
	"net/http"
)

type HttpErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func reportHttpError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(HttpErrorResponse{
		Error:  err.Error(),
		Status: statusCode,
	})
}
