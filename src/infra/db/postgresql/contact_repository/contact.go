package contactrepository

import (
	"github.com/rluisb/agenda-app/src/domain/usecases"
	"github.com/rluisb/agenda-app/src/domain/usecases/models"
	"github.com/rluisb/agenda-app/src/infra/db/postgresql/protocols"
)

type ContactPostgresRepository struct {
	Db protocols.Postgres
}

func NewContactPostgresRepository(db protocols.Postgres) *ContactPostgresRepository {
	return &ContactPostgresRepository{db}
}

func (c ContactPostgresRepository) Add(contact *usecases.AddContactModel) (*models.ContactModel, error) {
	contactModel := *models.NewContactModel("1", contact.Name, contact.Email, contact.Phone, contact.Address)
	_, err := c.Db.Exec("INSERT INTO contacts (name, email, phone, address) VALUES ($1, $2, $3, $4, $5)", contactModel.ID, contactModel.Name, contactModel.Email, contactModel.Phone, contactModel.Address)
	if err != nil {
		return nil, err
	}
	
	return &contactModel, nil
}