package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rluisb/agenda-app/src/domain/usecases"
	"github.com/rluisb/agenda-app/src/presentation/custom_errors"
	"github.com/rluisb/agenda-app/src/presentation/helpers"
	"github.com/rluisb/agenda-app/src/validation/protocols"
)

type CreateContactController struct{
	EmailValidator protocols.EmailValidator
	AddContact usecases.AddContact
}

func NewCreateContactController(emailValidator protocols.EmailValidator, addContact usecases.AddContact) *CreateContactController {
	return &CreateContactController{emailValidator, addContact}
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
	isValid, err := c.EmailValidator.IsValid(body["Email"])
	if err != nil {
		errorMessage := custom_errors.NewInternalServerError().Build()
		httpResponse := helpers.InternalServerError(errorMessage)
		w.WriteHeader(httpResponse.StatusCode)
		json.NewEncoder(w).Encode(httpResponse.Body)
		return
	}
	if !isValid {
		errorMessage := custom_errors.NewInvalidParamError(strings.ToLower("Email")).Build()
		httpResponse := helpers.BadRequest(errorMessage)
		w.WriteHeader(httpResponse.StatusCode)
		json.NewEncoder(w).Encode(httpResponse.Body)
		return
	}

	newContact, err := c.AddContact.Add(usecases.NewAddContactModel(body["Name"], body["Email"], body["Phone"], body["Address"]))
	if err != nil {
		errorMessage := custom_errors.NewInternalServerError().Build()
		httpResponse := helpers.InternalServerError(errorMessage)
		w.WriteHeader(httpResponse.StatusCode)
		json.NewEncoder(w).Encode(httpResponse.Body)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newContact)
}
