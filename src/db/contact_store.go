package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/rluisb/agenda-app/src/helper"
	"github.com/rluisb/agenda-app/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var (
	DBNAME = os.Getenv("MONGODB_DBNAME")
)

const COLLECTION = "contacts"

type ContactStore interface {
	GetContactByID(context.Context, string) (*types.Contact, error)
	GetContacts(context.Context, *types.GetContactsListQueryParams) ([]*types.Contact, error)
	CreateContact(context.Context, *types.Contact) (*types.Contact, error)
	DeleteContact(context.Context, string) error
	UpdateContact(context.Context, string, *types.UpdateContactParams) error
	PatchContact(context.Context, string, *types.UpdateContactParams) error
	Exists(context.Context, bson.M) (bool, error)
}

type MongoContactStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoContactStore(client *mongo.Client) *MongoContactStore {
	collection  := client.Database(DBNAME).Collection(COLLECTION)
	namePhoneEmailIndex := mongo.IndexModel{
    Keys: bson.D{
        {"name", 1},
        {"phone", 1},
				{"email", 1},
    }, Options: options.Index().SetUnique(true),
}
	name, err := collection.Indexes().CreateOne(context.Background(), namePhoneEmailIndex)
	if err != nil {
		panic(err)
	}
	log.Printf("Created index %s", name)

	return &MongoContactStore{
		client: client,
		coll: collection,
	}
}

func (s *MongoContactStore) GetContacts(ctx context.Context, queryParams *types.GetContactsListQueryParams) ([]*types.Contact, error) {
	var contacts []*types.Contact

	deletedAtQuery := bson.M{"$eq": 0}
	if(!queryParams.Active) {
		deletedAtQuery = bson.M{"$ne": 0}		
	}

	query := bson.D{{
			Key: "name", Value: primitive.Regex{Pattern: queryParams.Name, Options: "i"},
		}, {
			Key: "deleted_at", Value: deletedAtQuery,
	}}

	if queryParams.Name == "" {
		query = bson.D{{
			Key: "deleted_at", Value: deletedAtQuery,
		}}
	}

	queryOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	

	cursor, err := s.coll.Find(ctx, query, queryOptions)
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
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("contact not found")
		}
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
		if mongo.IsDuplicateKeyError(err) {
			errorMessage := helper.ParseMongoError(err.Error())
			return nil, errors.New(errorMessage)
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

func (s *MongoContactStore) UpdateContact(ctx context.Context, id string, updateContactParams *types.UpdateContactParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"name": updateContactParams.Name,
			"phone": updateContactParams.Phone,
			"email": updateContactParams.Email,
			"address": updateContactParams.Address,
			"updated_at": time.Now().Unix(),
		},
	} 
	exists, err := s.Exists(ctx, bson.M{
		"name": updateContactParams.Name,
		"phone": updateContactParams.Phone,
		"email": updateContactParams.Email,
	})
	if err != nil {
		return err
	}
	if exists {
		return errors.New("contact already exists")
	}
	result, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			errorMessage := helper.ParseMongoError(err.Error())
			return errors.New(errorMessage)
		}
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("contact not found")
	}
	return nil
}

func (s *MongoContactStore) PatchContact(ctx context.Context, id string, updateContactParams *types.UpdateContactParams) error {	
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now().Unix(),
		},
	} 
	if updateContactParams.Name != "" {
		update["$set"].(bson.M)["name"] = updateContactParams.Name
	}
	if updateContactParams.Phone != "" {
		update["$set"].(bson.M)["phone"] = updateContactParams.Phone
	}
	if updateContactParams.Email != "" {
		update["$set"].(bson.M)["email"] = updateContactParams.Email
	}
	if updateContactParams.Address != "" {
		update["$set"].(bson.M)["address"] = updateContactParams.Address
	}
	exists, err := s.Exists(ctx, bson.M{
		"name": updateContactParams.Name,
		"phone": updateContactParams.Phone,
		"email": updateContactParams.Email,
	})
	if err != nil {
		return err
	}
	if exists {
		return errors.New("contact already exists")
	}
	result, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			errorMessage := helper.ParseMongoError(err.Error())
			return errors.New(errorMessage)
		}
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("contact not found")
	}
	return nil
}

func (s *MongoContactStore) Exists(ctx context.Context, query bson.M) (bool, error) {
	count, err := s.coll.CountDocuments(ctx, query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

