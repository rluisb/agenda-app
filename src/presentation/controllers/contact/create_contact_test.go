package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rluisb/agenda-app/src/domain/usecases"
	"github.com/rluisb/agenda-app/src/domain/usecases/models"
)

var IsValid func(email string) (bool, error)

type GenericSpy struct {
	CallCount  int
	CalledWith interface{}
}

func NewGenericSpy() *GenericSpy {
	return &GenericSpy{}
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
	addContactStub := NewAddContactStub()
	emailValidatorStub := NewEmailValidator()
	return NewCreateContactController(emailValidatorStub, addContactStub)
}

func TestCreateContactBadRequest_MissingRequiredField(t *testing.T) {
	IsValid = func(email string) (bool, error) {
		return true, nil
	}
	sut := makeSut()

	
	table := map[string]usecases.AddContactModel{
		"name": *usecases.NewAddContactModel("", "john.doe@mail.com", "1234567890", "123 Main St"),
		"email":  *usecases.NewAddContactModel("john.doe", "", "1234567890", "123 Main St"),
		"phone": *usecases.NewAddContactModel("john.doe","john.doe@mail.com","","123 Main St"),
		"address": *usecases.NewAddContactModel("john.doe", "john.doe@mail.com", "1234567890", ""),
	}

	for missingField, contact := range table {
		t.Run("Should return badRequest when " + missingField + " is missing", func(t *testing.T) {
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
		})
	}
}

func TestCreateContactBadRequest_InvalidEmailProvided(t *testing.T) {
	emailValidatorSpy := NewGenericSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.CallCount++
		emailValidatorSpy.CalledWith = email
		return false, nil
	}
	sut := makeSut()
	contact := usecases.NewAddContactModel("John Doe", "invalid_email@mail.com", "1234567890", "123 Main St")
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
		if emailValidatorSpy.CallCount != 1 {
			t.Errorf("Expected email validator to be called 1 time, got %v", emailValidatorSpy.CallCount)
		}
		if emailValidatorSpy.CalledWith != contact.Email {
			t.Errorf("Expected email validator to have been called with %v, got %v", contact.Email, emailValidatorSpy.CalledWith)
		}
	}
}

func TestCreateContactInternalServerError_IfEmailValidatorFails(t *testing.T) {
	emailValidatorSpy := NewGenericSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.CallCount++
		emailValidatorSpy.CalledWith = email
		return false, errors.New("something went wrong")
	}
	sut := makeSut()
	contact := usecases.NewAddContactModel("John Doe", "invalid_email@mail.com", "1234567890", "123 Main St")
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
		if emailValidatorSpy.CallCount != 1 {
			t.Errorf("Expected email validator to be called 1 time, got %v", emailValidatorSpy.CallCount)
		}
		if emailValidatorSpy.CalledWith != contact.Email {
			t.Errorf("Expected email validator to have been called with %v, got %v", contact.Email, emailValidatorSpy.CalledWith)
		}
	}
}

var add func(contact *usecases.AddContactModel) (*models.ContactModel, error)

type AddContactStub struct {
}

func NewAddContactStub() AddContactStub {
	return AddContactStub{}
}

func (a AddContactStub) Add(contact *usecases.AddContactModel) (*models.ContactModel, error) {
	return add(contact)
}

func TestCreateContactWithSuccess(t *testing.T) {
	addAccountSpy := NewGenericSpy()
	add = func(contact *usecases.AddContactModel) (*models.ContactModel, error){
		addAccountSpy.CallCount++
		addAccountSpy.CalledWith = contact
		return &models.ContactModel{
			ID:      "1",
			Name:    contact.Name,
			Email:   contact.Email,
			Phone:   contact.Phone,
			Address: contact.Address,
		}, nil
	}
	emailValidatorSpy := NewGenericSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.CallCount++
		emailValidatorSpy.CalledWith = email
		return true, nil
	}
	sut := makeSut()
	contact := usecases.NewAddContactModel("John Doe", "john.doe@mail.com", "1234567890", "123 Main St")
	body, _ := json.Marshal(contact)
	r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	sut.handle(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v with body", http.StatusCreated, w.Code)
	}
	if w.Code == http.StatusCreated {
		expectedResponse := map[string]string{}
		var responseBody map[string]string
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		if len(responseBody) == 0 &&  len(expectedResponse) == 0 && len(responseBody) == len(expectedResponse) != true {
			t.Errorf("Expected response body %v, got %v", expectedResponse, expectedResponse)
		}
		if w.Header()["Content-Type"][0] != "application/json" {
			t.Errorf("Expected response header Content-Type to be application/json, got %v", w.HeaderMap["Content-Type"][0])
		}
		if w.Header()["Location"][0] != "http://localhost:8080/contacts/1" {
			t.Errorf("Expected response header Location to be http://localhost:8080/contacts/1, got %v", w.HeaderMap["Location"])
		}
		if emailValidatorSpy.CallCount != 1 {
			t.Errorf("Expected email validator to have been called 1 time, got %v", emailValidatorSpy.CallCount)
		}
		if reflect.DeepEqual(contact.Email, emailValidatorSpy.CalledWith) != true {
			t.Errorf("Expected email validator to have been called with %v, got %v", contact.Email, emailValidatorSpy.CalledWith)
		}
		if addAccountSpy.CallCount != 1 {
			t.Errorf("Expected add account to have been called 1 time, got %v", addAccountSpy.CallCount)
		}
		if reflect.DeepEqual(contact, addAccountSpy.CalledWith) != true {
			t.Errorf("Expected add account to been called with %v, got %v", contact, addAccountSpy.CalledWith)
		}
	}
}

func TestCreateContactInternalServerError_IfAccContactFails(t *testing.T) {
	addAccountSpy := NewGenericSpy()
	add = func(contact *usecases.AddContactModel) (*models.ContactModel, error) {
		addAccountSpy.CallCount++
		addAccountSpy.CalledWith = contact
		return nil, errors.New("something went wrong")
	}
	emailValidatorSpy := NewGenericSpy()
	IsValid = func(email string) (bool, error) {
		emailValidatorSpy.CallCount++
		emailValidatorSpy.CalledWith = email
		return true, nil
	}
	sut := makeSut()
	contact := usecases.NewAddContactModel("John Doe", "invalid_email@mail.com", "1234567890", "123 Main St")
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
		if emailValidatorSpy.CallCount != 1 {
			t.Errorf("Expected email validator to have been called 1 time, got %v", emailValidatorSpy.CallCount)
		}
		if reflect.DeepEqual(contact.Email, emailValidatorSpy.CalledWith) != true {
			t.Errorf("Expected email validator to have been called with %v, got %v", contact.Email, emailValidatorSpy.CalledWith)
		}
		if addAccountSpy.CallCount != 1 {
			t.Errorf("Expected add account to have been called 1 time, got %v", addAccountSpy.CallCount)
		}
		if reflect.DeepEqual(contact, addAccountSpy.CalledWith) != true {
			t.Errorf("Expected add account to been called with %v, got %v", contact, addAccountSpy.CalledWith)
		}
	}
}
