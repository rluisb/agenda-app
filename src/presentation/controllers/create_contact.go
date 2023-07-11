package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rluisb/agenda-app/src/presentation/custom_errors"
	"github.com/rluisb/agenda-app/src/presentation/helpers"
	"github.com/rluisb/agenda-app/src/presentation/protocols"
)

type CreateContactController struct{
	EmailValidator protocols.EmailValidator
}

func NewCreateContactController(emailValidator protocols.EmailValidator) *CreateContactController {
	return &CreateContactController{emailValidator}
}

func (c CreateContactController) handle(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	requiredFields := []string{"Name", "Email", "Phone", "Address"}
	json.NewDecoder(r.Body).Decode(&body)
	for _, field := range requiredFields {
		if body[field] == "" {
			errorMessage := custom_errors.NewMissingParamError(field).Build()
			httpResponse := helpers.BadRequest(errorMessage)
			w.WriteHeader(httpResponse.StatusCode)
			json.NewEncoder(w).Encode(httpResponse.Body)
			return
		}
	}
	if !c.EmailValidator.IsValid(body["Email"]) {
		errorMessage := custom_errors.NewInvalidParamError(strings.ToLower("Email")).Build()
		httpResponse := helpers.BadRequest(errorMessage)
		w.WriteHeader(httpResponse.StatusCode)
		json.NewEncoder(w).Encode(httpResponse.Body)
		return
	}
}
