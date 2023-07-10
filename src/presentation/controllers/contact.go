package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/rluisb/agenda-app/src/presentation/helpers"
)

type ContactController struct{}

func NewContactController() *ContactController {
	return &ContactController{}
}

func (c ContactController) handle(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	if body["name"] == "" {
		httpResponse := helpers.BadRequest("name")
		w.WriteHeader(httpResponse.StatusCode)
		json.NewEncoder(w).Encode(httpResponse.Body)
	}
	if body["email"] == "" {
		httpResponse := helpers.BadRequest("email")
		w.WriteHeader(httpResponse.StatusCode)
		json.NewEncoder(w).Encode(httpResponse.Body)
	}
}
