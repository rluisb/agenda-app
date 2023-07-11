package utils

type EmailValidatorAdapter struct {}

func NewEmailValidatorAdapter() *EmailValidatorAdapter {
	return &EmailValidatorAdapter{}
}

func (v EmailValidatorAdapter) IsValid(email string) (bool, error) {
	return false, nil
}