package repository

import (
	"context"
	"log"
	"logger/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository interface {
	Insert(entry models.LogEntry) error
	All() ([]*models.LogEntry, error)
	GetOne(id string) (*models.LogEntry, error)
	DropCollection() error
	Update(payload models.LogEntry) (*mongo.UpdateResult, error)
}

type logRepository struct {
	client    *mongo.Client
	dbTimeout time.Duration
}

func NewLogRepository(client *mongo.Client, dbTimeout time.Duration) LogRepository {
	return &logRepository{
		client:    client,
		dbTimeout: dbTimeout,
	}
}

func (repo *logRepository) Insert(entry models.LogEntry) error {
	collection := repo.client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil

}

func (repo *logRepository) All() ([]*models.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.dbTimeout)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.LogEntry

	for cursor.Next(ctx) {
		var item models.LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil

}

func (repo *logRepository) GetOne(id string) (*models.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.dbTimeout)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry models.LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil

}
func (repo *logRepository) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil

}

func (repo *logRepository) Update(payload models.LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := repo.client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", payload.Name},
				{"data", payload.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
