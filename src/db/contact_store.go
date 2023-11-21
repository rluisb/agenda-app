package db

import (
	"context"
	"errors"

	"github.com/rluisb/agenda-app/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const contactColl = "contacts"

type ContactStore interface {
	GetContactByID(context.Context, string) (*types.Contact, error)
	GetContacts(context.Context) ([]*types.Contact, error)
	CreateContact(context.Context, *types.Contact) (*types.Contact, error)
}

type MongoContactStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoContactStore(client *mongo.Client) *MongoContactStore {
	return &MongoContactStore{
		client: client,
		coll: client.Database(DBNAME).Collection(contactColl),
	}
}

func (s *MongoContactStore) GetContacts(ctx context.Context) ([]*types.Contact, error) {
	var contacts []*types.Contact
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if cursor.RemainingBatchLength() == 0 {
		return []*types.Contact{}, nil
	}
	if err := cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *MongoContactStore) GetContactByID(ctx context.Context, id string) (*types.Contact, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	var contact types.Contact
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *MongoContactStore) CreateContact(ctx context.Context, contact *types.Contact) (*types.Contact, error) {
	res, err := s.coll.InsertOne(ctx, contact)
	if err != nil {
		return nil, err
	}
	contact.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return contact, nil
}



