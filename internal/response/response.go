package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// DATA STRUCTURE CONCEPT:
// APIResponse is a structured wrapper.
// Think of it as a fixed schema container — like a struct-based map
// enforcing consistent response shape across system.

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}
