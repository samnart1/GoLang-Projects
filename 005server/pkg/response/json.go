package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// if we can't encode, send internal server error and return
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"message": 	message,
			"status":	statusCode,
		},
	}

	JSON(w, statusCode, errorResponse)
}

func Success(w http.ResponseWriter, data interface{}) {
	response := map[string]interface{}{
		"success":	true,
		"data":		data,
	}

	JSON(w, http.StatusOK, response)
}