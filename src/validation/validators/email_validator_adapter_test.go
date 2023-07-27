package utils

import (
	"errors"
	"net/mail"
	"reflect"
	"testing"
)

type GenericSpy struct {
	CallCount int
	CalledWith interface{}
}

var Parse func (string) (*mail.Address, error)

type AddressParserStub struct {}

func NewAddressParserStub() *AddressParserStub {
	return &AddressParserStub{}
}

func (a AddressParserStub) Parse(email string) (*mail.Address, error) {
	return Parse(email)
}

func TestEmailValidator_ReturnFalse(t *testing.T) {
	addressParserSpy := GenericSpy{}
	Parse = func (email string) (*mail.Address, error) {
		addressParserSpy.CallCount++
		addressParserSpy.CalledWith = email
		return nil, nil
	}
	addressParserStub := NewAddressParserStub()
	sut := NewEmailValidatorAdapter(addressParserStub)
	email := "invalid_email@mail.com"
	isValid, err := sut.IsValid(email)
	if isValid != false{
		t.Errorf("Expected email validator to return false, got %v", isValid)
	}
	if isValid == false && err != nil {
		t.Errorf("Expected email validator to return false with no error, got %v and error %v", isValid, err)
	}
	if addressParserSpy.CallCount != 1 {
		t.Errorf("Expected address parser to be called once, got %v", addressParserSpy.CallCount)
	}
	if !reflect.DeepEqual(addressParserSpy.CalledWith, email) {
		t.Errorf("Expected address parser to be called with %v, got %v", email, addressParserSpy.CalledWith)
	}
}

func TestEmailValidator_ReturnTrue(t *testing.T) {
	addressParserSpy := GenericSpy{}
	Parse = func (email string) (*mail.Address, error) {
		addressParserSpy.CallCount++
		addressParserSpy.CalledWith = email
		return &mail.Address{}, nil
	}
	addressParserStub := NewAddressParserStub()
	sut := NewEmailValidatorAdapter(addressParserStub)
	email := "valid_email@mail.com"
	isValid, _ := sut.IsValid(email)
	if isValid != true{
		t.Errorf("Expected email validator to return true, got %v", isValid)
	}
	if addressParserSpy.CallCount != 1 {
		t.Errorf("Expected address parser to be called once, got %v", addressParserSpy.CallCount)
	}
	if !reflect.DeepEqual(addressParserSpy.CalledWith, email) {
		t.Errorf("Expected address parser to be called with %v, got %v", email, addressParserSpy.CalledWith)
	}
}


func TestEmailValidator_ReturnError(t *testing.T) {
	addressParserSpy := GenericSpy{}
	Parse = func (email string) (*mail.Address, error) {
		addressParserSpy.CallCount++
		addressParserSpy.CalledWith = email
		return nil, errors.New("something went wrong")
	}
	addressParserStub := NewAddressParserStub()
	sut := NewEmailValidatorAdapter(addressParserStub)
	email := "valid_email@mail.com"
	isValid, err:= sut.IsValid(email)
	if isValid != false{
		t.Errorf("Expected email validator to return false, got %v", isValid)
	}
	if addressParserSpy.CallCount != 1 {
		t.Errorf("Expected address parser to be called once, got %v", addressParserSpy.CallCount)
	}
	if !reflect.DeepEqual(addressParserSpy.CalledWith, email) {
		t.Errorf("Expected address parser to be called with %v, got %v", email, addressParserSpy.CalledWith)
	}
	if err == nil {
		t.Errorf("Expected email validator to return error, got %v", err)
	}
}