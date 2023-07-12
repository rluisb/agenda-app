package protocols

import (
	"github.com/rluisb/agenda-app/src/domain/usecases"
	"github.com/rluisb/agenda-app/src/domain/usecases/models"
)

type AddContactRepository interface {
	Add(contact *usecases.AddContactModel) (*models.ContactModel, error)
}