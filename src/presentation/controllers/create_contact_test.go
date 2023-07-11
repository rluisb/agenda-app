package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var IsValid func(email string) (bool, error)

type EmailValidatorSpy struct {
	IsValidCallCount  int
	IsValidCalledWith string
}

func NewEmailValidatorSpy() *EmailValidatorSpy {
	return &EmailValidatorSpy{}
}

func (spy *EmailValidatorSpy) reset() {
	spy.IsValidCallCount = 0
	spy.IsValidCalledWith = ""
}

type EmailValidatorStub struct {
}

func NewEmailValidator() EmailValidatorStub {
	return EmailValidatorStub{}
}

func (e EmailValidatorStub) IsValid(email string) (bool, error) {
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
	IsValid = func(email string) (bool, error) {
		return true, nil
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
	emailValidatorSpy := NewEmailValidatorSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.IsValidCallCount++
		emailValidatorSpy.IsValidCalledWith = email
		return false, nil
	}
	sut := makeSut()
	contact := Contact{
		Name:    "John Doe",
		Email:   "invalid_email@mail.com",
		Phone:   "1234567890",
		Address: "123 Main St",
	}
	body, _ := json.Marshal(contact)
	r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	sut.handle(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v with body", http.StatusBadRequest, w.Code)
	}
	if w.Code == http.StatusBadRequest {
		expectedResponse := map[string]string{"message": "invalid email"}
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		if reflect.DeepEqual(expectedResponse, responseBody) != true {
			t.Errorf("Expected response body %v, got %v", expectedResponse, responseBody)
		}
		if emailValidatorSpy.IsValidCallCount != 1 {
			t.Errorf("Expected email validator to be called 1 time, got %v", emailValidatorSpy.IsValidCallCount)
		}
		if emailValidatorSpy.IsValidCalledWith != contact.Email {
			t.Errorf("Expected email validator to be called with %v, got %v", contact.Email, emailValidatorSpy.IsValidCalledWith)
		}
	}
}

func TestCreateContactWithSuccess_CorrectEmail(t *testing.T) {
	emailValidatorSpy := NewEmailValidatorSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.IsValidCallCount++
		emailValidatorSpy.IsValidCalledWith = email
		return true, nil
	}
	sut := makeSut()
	contact := Contact{
		Name:    "John Doe",
		Email:   "john.doe@mail.com",
		Phone:   "1234567890",
		Address: "123 Main St",
	}
	body, _ := json.Marshal(contact)
	r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	sut.handle(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v with body", http.StatusOK, w.Code)
	}
	if w.Code == http.StatusOK {
		expectedResponse := map[string]string{}
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		if len(responseBody) != len(expectedResponse) {
			t.Errorf("Expected response body %v, got %v", expectedResponse, responseBody)
		}
	}
}

func TestCreateContactInternalServerError_IfEmailValidatorThrows(t *testing.T) {
	emailValidatorSpy := NewEmailValidatorSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.IsValidCallCount++
		emailValidatorSpy.IsValidCalledWith = email
		return false, errors.New("something went wrong")
	}
	sut := makeSut()
	contact := Contact{
		Name:    "John Doe",
		Email:   "invalid_email@mail.com",
		Phone:   "1234567890",
		Address: "123 Main St",
	}
	body, _ := json.Marshal(contact)
	r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	sut.handle(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v with body", http.StatusBadRequest, w.Code)
	}
	if w.Code == http.StatusInternalServerError {
		expectedResponse := map[string]string{"message": "internal server error"}
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		if reflect.DeepEqual(expectedResponse, responseBody) != true {
			t.Errorf("Expected response body %v, got %v", expectedResponse, responseBody)
		}
		if emailValidatorSpy.IsValidCallCount != 1 {
			t.Errorf("Expected email validator to be called 1 time, got %v", emailValidatorSpy.IsValidCallCount)
		}
		if emailValidatorSpy.IsValidCalledWith != contact.Email {
			t.Errorf("Expected email validator to be called with %v, got %v", contact.Email, emailValidatorSpy.IsValidCalledWith)
		}
	}
}