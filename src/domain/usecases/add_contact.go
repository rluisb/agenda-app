package usecases

import "github.com/rluisb/agenda-app/src/domain/usecases/models"

type AddContactModel struct {
	Name    string
	Email   string
	Phone   string
	Address string
}

func NewAddContactModel(name string, email string, phone string, address string) *AddContactModel {
	return &AddContactModel{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
	}
}

type AddContact interface {
	Add(contact *AddContactModel) (*models.ContactModel, error)
}