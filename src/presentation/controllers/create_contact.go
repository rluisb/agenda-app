package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/rluisb/agenda-app/src/presentation/helpers"
)

type CreateContactController struct{}

func NewCreateContactController() *CreateContactController {
	return &CreateContactController{}
}

func (c CreateContactController) handle(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	requiredFields := []string{"Name", "Email", "Phone", "Address"}
	json.NewDecoder(r.Body).Decode(&body)
	for _, field := range requiredFields {
		if body[field] == "" {
			httpResponse := helpers.BadRequest(field)
			w.WriteHeader(httpResponse.StatusCode)
			json.NewEncoder(w).Encode(httpResponse.Body)
			return
		}
	}
}
