package contactrepository

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
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

var Exec func(query string, args ...interface{}) (sql.Result, error)

type PostgrsMock struct {
	DB *sql.DB
}

func NewPostgresStub(db *sql.DB) *PostgrsMock {
	return &PostgrsMock{db}
}

func (p PostgrsMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	return Exec(query, args...)
}

var LastInsertId func() (int64, error)
var RowsAffected func() (int64, error)

type SqlResultStub struct{}

func NewSqlResultStub() *SqlResultStub {
	return &SqlResultStub{}
}

func (s SqlResultStub) LastInsertId() (int64, error) {
	return LastInsertId()
}

func (s SqlResultStub) RowsAffected() (int64, error) {
	return RowsAffected()
}

func TestContactPostgresRepository(t *testing.T) {
	t.Run("Should return a contact model with an ID", func(t *testing.T) {
		postgresExecSpy := NewGenericSpy()
		Exec = func(query string, args ...interface{}) (sql.Result, error) {
			postgresExecSpy.CallCount++
			postgresExecSpy.CalledWith = query
			sqlResultStub := NewSqlResultStub()
			return sqlResultStub, nil
		}
		
		postgresStub := NewPostgresStub(&sql.DB{})
		sut := NewContactPostgresRepository(postgresStub)
		addContactModel := usecases.NewAddContactModel("John Doe", "john.doe@mail.com", "1234567890", "123 Main St")
		expectedContactModel := models.NewContactModel("1", addContactModel.Name, addContactModel.Email, addContactModel.Phone, addContactModel.Address)
		contact, err := sut.Add(addContactModel)
		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}
		if contact == nil {
			t.Errorf("Expected contact to be not nil, got %v", contact)
		}
		if contact.ID != "1" {
			t.Errorf("Expected contact ID to be 1, got %v", contact.ID)
		}
		if reflect.DeepEqual(expectedContactModel, contact) != true{
			t.Errorf("Expected contact to be %v, got %v", expectedContactModel, contact)
		}
	}) 
}