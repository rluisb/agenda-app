package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var IsValid func (email string) bool

type EmailValidatorStub struct {
}

func NewEmailValidator() EmailValidatorStub {
	return EmailValidatorStub{}
}

func (e EmailValidatorStub) IsValid(email string) bool { 
	return IsValid(email)
}

func makeSut() *CreateContactController {
	emailValidatorStub := NewEmailValidator()
	return NewCreateContactController(emailValidatorStub)
}

type Contact struct {
	Name    string
	Email   string
	Phone   string
	Address string
}

func TestCreateContactBadRequest_MissingRequiredField(t *testing.T) {
	IsValid = func (email string) bool {
		return true
	}

	sut := makeSut()

	table := map[string]Contact{
		"name": {
			Email:   "john.doe@mail.com",
			Phone:   "1234567890",
			Address: "123 Main St",
		},
		"email": {
			Name:    "john.doe",
			Phone:   "1234567890",
			Address: "123 Main St",
		},
		"phone": {
			Name:    "john.doe",
			Email:   "john.doe@mail.com",
			Address: "123 Main St",
		},
		"address": {
			Name:  "john.doe",
			Email: "john.doe@mail.com",
			Phone: "1234567890",
		},
	}

	for missingField, contact := range table {
		body, _ := json.Marshal(contact)
		r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		sut.handle(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v with body", http.StatusBadRequest, w.Code)
		}
		if w.Code == http.StatusBadRequest {
			expected := map[string]string{"message": missingField + " is required"}
			var responseBody map[string]string
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			if reflect.DeepEqual(expected, responseBody) != true {
				t.Errorf("Expected response body %v, got %v", expected, responseBody)
			}
		}
	}
}

func TestCreateContactBadRequest_InvalidEmailProvided(t *testing.T) {
	IsValid = func (email string) bool {
		return false
	}
	sut := makeSut()
	contatc := Contact{
		Name:    "John Doe",
		Email:   "invalid_email@mail.com",
		Phone:   "1234567890",
		Address: "123 Main St",
	}
	body, _ := json.Marshal(contatc)
	r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	sut.handle(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v with body", http.StatusBadRequest, w.Code)
	}
	if w.Code == http.StatusBadRequest {
		expected := map[string]string{"message": "invalid email"}
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		if reflect.DeepEqual(expected, responseBody) != true {
			t.Errorf("Expected response body %v, got %v", expected, responseBody)
		}
	}
}
