package utils

import "net/mail"

type AddressParser interface {
	Parse(email string) (*mail.Address, error)
}

type EmailValidatorAdapter struct {
	AddressParser AddressParser
}

func NewEmailValidatorAdapter(addressParser AddressParser) *EmailValidatorAdapter {
	return &EmailValidatorAdapter{addressParser}
}

func (v EmailValidatorAdapter) IsValid(email string) (bool, error) {
	isValid, err := v.AddressParser.Parse(email)
	if err != nil {
		return false, err
	}
	return isValid != nil, nil
}