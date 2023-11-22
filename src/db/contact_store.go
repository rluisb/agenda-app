package db

import (
	"context"
	"errors"
	"time"

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
	DeleteContact(context.Context, string) error
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
	currentTime := time.Now().Unix()
	contact.CreatedAt = currentTime
	contact.UpdatedAt = currentTime
	res, err := s.coll.InsertOne(ctx, contact)
	if err != nil {
		return nil, err
	}
	contact.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return contact, nil
}

func (s *MongoContactStore) DeleteContact(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	currentTime := time.Now().Unix()
	_, err = s.coll.UpdateByID(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"updated_at": currentTime, "deleted_at": currentTime}})
	if err != nil {
		return err
	}
	return nil
}



