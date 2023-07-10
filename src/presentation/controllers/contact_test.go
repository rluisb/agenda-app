package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateContactBadRequest_NoNameProvided(t *testing.T) {
	sut := NewContactController()
	body, _ := json.Marshal(map[string]string{
		"email":   "john.doe@mail.com",
		"phone":   "1234567890",
		"address": "123 Main St",
	})
	r, _ := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	sut.handle(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, w.Code)
	}
}
