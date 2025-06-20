package db

import (
	"context"
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
    findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}) // ascending

    cursor, err := collection.Find(context.TODO(), filter, findOptions)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var msg ChatMessage
        if err := cursor.Decode(&msg); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}