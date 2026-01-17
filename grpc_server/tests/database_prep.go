package tests

import (
	"context"
	running_trackerpb "grpcserver/proto/generated_files"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainerDocument struct {
	ID               primitive.ObjectID `bson:"_id"`
	Brand            string             `bson:"brand"`
	Model            string             `bson:"model"`
	PurchaseDate     string             `bson:"purchase_date"`
	ExpectedLifespan int32              `bson:"expected_lifespan"`
	SurfaceType      int32              `bson:"surface_type"`
	Status           int32              `bson:"status"`
}

func insertTestTrainer(t *testing.T, client *mongo.Client, brand, model string) TrainerDocument {
	ctx := context.Background()
	coll := client.Database("data").Collection("trainers")

	trainer := TrainerDocument{
		ID:               primitive.NewObjectID(),
		Brand:            brand,
		Model:            model,
		PurchaseDate:     "2024-01-01",
		ExpectedLifespan: 700,
		SurfaceType:      int32(running_trackerpb.SurfaceType_ROAD),
		Status:           int32(running_trackerpb.TrainerStatus_NEW),
	}

	_, err := coll.InsertOne(ctx, trainer)
	if err != nil {
		t.Fatalf("Failed to insert trainer: %v", err)
	}
	return trainer
}
