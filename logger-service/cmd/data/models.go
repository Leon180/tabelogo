package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Models struct {
	LogEntry LogEntry
}

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) InsertOne(entry LogEntry) error {
	// get collection logs
	collection := client.Database("logs").Collection("logs")
	// insert one log entry
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry", err)
		return err
	}
	return nil
}

func (l *LogEntry) FindAll() ([]*LogEntry, error) {
	// context to set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// get collection logs
	collection := client.Database("logs").Collection("logs")
	// find all log entries and sort by created_at in descending order
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error getting log entries", err)
		return nil, err
	}
	// close cursor after function ends with timeout context (the context will be canceled when the function ends or times out)
	defer cursor.Close(ctx)

	var entries []*LogEntry
	for cursor.Next(ctx) {
		var entry LogEntry
		if err := cursor.Decode(&entry); err != nil {
			log.Println("Error decoding log entry", err)
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (l *LogEntry) FindByID(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting string to objectID", err)
		return nil, err
	}

	var entry LogEntry
	if err := collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err != nil {
		log.Println("Error getting log entry", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Updata() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error converting string to objectID", err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)
	if err != nil {
		log.Println("Error updating log entry", err)
		return nil, err
	}

	return result, nil
}
