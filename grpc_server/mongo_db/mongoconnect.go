package mongodb

import (
	"context"
	"fmt"
	"grpcserver/internals/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")

	if username == "" || password == "" {
		log.Fatal("MONGODB_USERNAME or MONGODB_PASSWORD not set")
	}

	ctx := context.Background()
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", username, password, host, port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to connect to database")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to ping to database")
	}

	utils.InfoLogger.Println("Connected to MongoDB")
	return client, nil
}

func DisconnectMongoClient(client *mongo.Client, ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		utils.ErrorLogger.Println("Failed to disconnect the database")
	}
}
