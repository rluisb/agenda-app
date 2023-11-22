package types

import (
	"net/mail"
)

type CreateContactParams struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func (params CreateContactParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.Name == "" {
		errors["name"] = "name is required"
	}
	if params.Phone == "" {
		errors["phone"] = "phone is required"
	}
	if params.Email == "" {
		errors["email"] = "email is required"
	}
	_, err := mail.ParseAddress(params.Email)
	if err != nil {
		errors["email"] = "email is invalid"
	}
	if params.Address == "" {
		errors["address"] = "address is required"
	}
	return errors
}

type Contact struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	Name			string `bson:"name" json:"name"`
	Phone			string `bson:"phone" json:"phone"`
	Email			string `bson:"email" json:"email"`
	Address		string `bson:"address" json:"address"`
	CreatedAt int64  `bson:"created_at" json:"-"`
	UpdatedAt int64  `bson:"updated_at" json:"-"`
	DeletedAt int64  `bson:"deleted_at" json:"-"`
}

func NewContactFromParams(params CreateContactParams) *Contact {
	return &Contact{
		Name: params.Name,
		Phone: params.Phone,
		Email: params.Email,
		Address: params.Address,
	}
}