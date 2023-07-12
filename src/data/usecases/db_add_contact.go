package usecases

import (
	"github.com/rluisb/agenda-app/src/data/usecases/protocols"
	"github.com/rluisb/agenda-app/src/domain/usecases"
	"github.com/rluisb/agenda-app/src/domain/usecases/models"
)

type DbAddContact struct {
	AddContactRepository protocols.AddContactRepository
}

func NewDbAddContact(addContactRepository protocols.AddContactRepository) *DbAddContact {
	return &DbAddContact{addContactRepository}
}

func (d DbAddContact) Add(addContactModel *usecases.AddContactModel) (*models.ContactModel, error) {
	return d.AddContactRepository.Add(addContactModel)
}