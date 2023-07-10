package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateContactBadRequest_MissingRequiredField(t *testing.T) {
	sut := NewContactController()
	type Contact struct {
		Name         string
		Email        string
		Phone        string
		Address      string
	}


	table := map[string]Contact{
		"name": {
			Email:        "john.doe@mail.com",
			Phone:        "1234567890",
			Address:      "123 Main St",
		},
		"email": {
			Name:         "john.doe",
			Phone:        "1234567890",
			Address:      "123 Main St",
		},
		"phone": {
			Name:         "john.doe",
			Email:        "john.doe@mail.com",
			Address:      "123 Main St",
		},
		"address": {
			Name:         "john.doe",
			Email:        "john.doe@mail.com",
			Phone:        "1234567890",
		},
	}

	for missingField, contact := range table {
		body, _ := json.Marshal(contact)
		r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		sut.CreateContact(w, r)
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
