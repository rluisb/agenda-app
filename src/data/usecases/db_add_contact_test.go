package usecases

import (
	"errors"
	"testing"

	"github.com/rluisb/agenda-app/src/domain/usecases"
	"github.com/rluisb/agenda-app/src/domain/usecases/models"
)

type GenericSpy struct {
	CallCount  int
	CalledWith interface{}
}

func NewGenericSpy() *GenericSpy {
	return &GenericSpy{}
}

var Add func(addContactModel *usecases.AddContactModel) (*models.ContactModel, error)

type AddContactRepositoryStub struct{}

func NewAddContactRepositoryStub() *AddContactRepositoryStub {
	return &AddContactRepositoryStub{}
}

func (a AddContactRepositoryStub) Add(addContactModel *usecases.AddContactModel) (*models.ContactModel, error) {
	return Add(addContactModel)
}

func TestDbAddContact(t *testing.T) {
	t.Run("Should call AddContactRepository with correct values", func(t *testing.T) {
		addContactRepositorySpy := NewGenericSpy()
		Add = func(addContactModel *usecases.AddContactModel) (*models.ContactModel, error) {
			addContactRepositorySpy.CallCount++
			addContactRepositorySpy.CalledWith = addContactModel
			contactModel := models.NewContactModel("1", "John Doe", "john.doe@mail.com", "1234567890", "123 Main St")
			return contactModel, nil
		}
		addContactRepositoryStub := NewAddContactRepositoryStub()
		sut := NewDbAddContact(addContactRepositoryStub)
		addContactModel := usecases.NewAddContactModel("John Doe", "john.doe@mail.com", "1234567890", "123 Main St")
		sut.Add(addContactModel)
		if addContactRepositorySpy.CallCount != 1 {
			t.Errorf("Expected call count to be 1, but got %d", addContactRepositorySpy.CallCount)
		}
		if addContactRepositorySpy.CalledWith != addContactModel {
			t.Errorf("Expected called with to be %v, but got %v", addContactModel, addContactRepositorySpy.CalledWith)
		}
	})

	t.Run("Should return error if AddContactRepository return error", func(t *testing.T) {
		addContactRepositorySpy := NewGenericSpy()
		Add = func(addContactModel *usecases.AddContactModel) (*models.ContactModel, error) {
			addContactRepositorySpy.CallCount++
			addContactRepositorySpy.CalledWith = addContactModel
			return nil, errors.New("something went wrong")
		}
		addContactRepositoryStub := NewAddContactRepositoryStub()
		sut := NewDbAddContact(addContactRepositoryStub)
		addContactModel := usecases.NewAddContactModel("John Doe", "john.doe@mail.com", "1234567890", "123 Main St")
		_, err := sut.Add(addContactModel)
		if addContactRepositorySpy.CallCount != 1 {
			t.Errorf("Expected call count to be 1, but got %d", addContactRepositorySpy.CallCount)
		}
		if addContactRepositorySpy.CalledWith != addContactModel {
			t.Errorf("Expected called with to be %v, but got %v", addContactModel, addContactRepositorySpy.CalledWith)
		}
		if err == nil {
			t.Errorf("Expected error to be returned")
		}
	})
}
