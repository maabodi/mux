package model

import (
	"encoding/json"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, code int, payload interface{}) {

	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ResponseError(w http.ResponseWriter, code int, message interface{}) {
	ResponseJson(w, code, message)
}

type ResponseUser struct {
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Role        string      `json:"role"`
	UserProfile UserProfile `json:"user_profile"`
	Avatar      Avatar      `json:"avatar"`
}
