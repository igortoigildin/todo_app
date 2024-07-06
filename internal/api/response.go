package api

import (
	"encoding/json"
	"net/http"
)

// method for error handling in JSON
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	result := make(map[string]interface{})
	result["error"] = "Ошибка"
	json.NewEncoder(w).Encode(result)
}
