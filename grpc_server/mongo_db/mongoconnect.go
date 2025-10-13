package mongodb

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient() (*mongo.Client, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("MONGODB_USERNAME or MONGODB_PASSWORD not set")
	}

	ctx := context.Background()
	mongoURI := fmt.Sprintf("mongodb://%s:%s@localhost:27017/?authSource=admin", username, password)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}