package db

import (
	"context"
	"fmt"
	"go-ai/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatMessage struct {
	UserID    string    `bson:"user_id,omitempty" json:"-"`
	Role      string    `bson:"role" json:"role"`
	Content   string    `bson:"content" json:"content"`
	Timestamp time.Time `bson:"timestamp,omitempty" json:"-"`
}

var client *mongo.Client
var collection *mongo.Collection

// InitMongo connects to MongoDB using env variables and sets up the collection
func InitMongo() {

	uri := config.GetMongoURI()
	dbName := config.GetMongoDB()
	collName := config.GetMongoCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	collection = client.Database(dbName).Collection(collName)
	log.Println("✅ Connected to MongoDB")
}

// StoreMessage saves a chat message in MongoDB
func StoreMessage(userID string, msg ChatMessage) error {
	msg.UserID = userID
	msg.Timestamp = time.Now()
	_, err := collection.InsertOne(context.Background(), msg)
	return err
}

// GetMessages retrieves all messages for a user
func GetMessages(userID string) ([]ChatMessage, error) {
	var messages []ChatMessage

	filter := bson.M{"user_id": userID}
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Printf("[ERROR] Find() failed: %v", err)
		return nil, fmt.Errorf("Find() failed: %w", err)
	}
	defer func() {
		if err := cursor.Close(context.TODO()); err != nil {
			log.Printf("[WARN] Failed to close cursor: %v", err)
		}
	}()

	count := 0
	for cursor.Next(context.TODO()) {
		var msg ChatMessage
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("[ERROR] Decode() failed: %v", err)
			return nil, fmt.Errorf("Decode() failed: %w", err)
		}
		messages = append(messages, msg)
		count++
	}

	if err := cursor.Err(); err != nil {
		log.Printf("[ERROR] Cursor iteration error: %v", err)
		return nil, fmt.Errorf("Cursor iteration error: %w", err)
	}

	return messages, nil
}
