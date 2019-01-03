package server

import (
	"net/http"
	"encoding/json"
)

func Error(w http.ResponseWriter, code int, message string) {
	JsonResponse(w, code, map[string]string{"error": message})
}

func JsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
