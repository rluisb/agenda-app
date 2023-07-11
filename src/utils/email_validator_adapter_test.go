package utils

import "testing"

func TestEmailValidator_ReturnFalse(t *testing.T) {
	sut := NewEmailValidatorAdapter()
	isValid, _ := sut.IsValid("invalid_email@mail.com")
	if isValid != false {
		t.Errorf("Expected email validator to return false, got %v", isValid)
	}
}