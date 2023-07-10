package controllers

import (
	"encoding/json"
	"net/http"
)

type ContactController struct{}

func NewContactController() *ContactController {
	return &ContactController{}
}

func (c ContactController) handle(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	if body["name"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "name is required"})
	}
	if body["email"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "email is required"})
	}
}
