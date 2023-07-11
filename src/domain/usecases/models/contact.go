package models

type ContactModel struct {
	ID      string
	Name    string
	Email   string
	Phone   string
	Address string
}

func NewContactModel(id string, name string, email string, phone string, address string) *ContactModel {
	return &ContactModel{
		ID:      id,
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
	}
}