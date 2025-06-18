package db

import "time"

type ChatEntry struct {
	UserID    string    `bson:"user_id"`
	Role      string    `bson:"role"`
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"timestamp"`
}