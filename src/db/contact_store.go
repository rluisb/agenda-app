package db

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/rluisb/agenda-app/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const contactColl = "contacts"

type ContactStore interface {
	GetContactByID(context.Context, string) (*types.Contact, error)
	GetContacts(context.Context, string) ([]*types.Contact, error)
	CreateContact(context.Context, *types.Contact) (*types.Contact, error)
	DeleteContact(context.Context, string) error
}

type MongoContactStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoContactStore(client *mongo.Client) *MongoContactStore {
	collection  := client.Database(DBNAME).Collection(contactColl)
	indexModel := mongo.IndexModel{
    Keys: bson.D{
        {"name", 1},
        {"phone", 1},
				{"email", 1},
    }, Options: options.Index().SetUnique(true),
}
	name, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}
	log.Printf("Created index %s", name)
	return &MongoContactStore{
		client: client,
		coll: collection,
	}
}

func (s *MongoContactStore) GetContacts(ctx context.Context, name string) ([]*types.Contact, error) {
	var contacts []*types.Contact

	if name != "" {
		cursor, err := s.coll.Find(ctx, bson.D{{
			Key: "name", Value: primitive.Regex{Pattern: name, Options: "i"},
		}, {
			Key: "deleted_at", Value: 0,
		}}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))

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

	cursor, err := s.coll.Find(ctx, bson.M{
		"deleted_at": 0,
	})
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
		if strings.Contains(err.Error(), "E11000") {
			return nil, errors.New("contact with name, phone and email already exists")
		}
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



