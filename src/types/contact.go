package types

import (
	"errors"
	"net/mail"
)

type CreateContactParams struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func (params CreateContactParams) Validate() error {
	if params.Name == "" {
		return errors.New("name is required")
	}
	if params.Phone == "" {
		return errors.New("phone is required")
	}
	if params.Email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(params.Email)
	if err != nil {
		return errors.New("invalid email")
	}
	if params.Address == "" {
		return errors.New("address is required")
	}
	return nil
}



type Contact struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	Name			string `bson:"name" json:"name"`
	Phone			string `bson:"phone" json:"phone"`
	Email			string `bson:"email" json:"email"`
	Address		string `bson:"address" json:"address"`
}

func NewContactFromParams(params CreateContactParams) (*Contact, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return &Contact{
		Name: params.Name,
		Phone: params.Phone,
		Email: params.Email,
		Address: params.Address,
	}, nil
}