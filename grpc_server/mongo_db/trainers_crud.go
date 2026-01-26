package mongodb

import (
	"context"
	"fmt"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	pb "grpcserver/proto/generated_files"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func AddTrainersToDB(ctx context.Context, request_trainers []*pb.Trainer) ([]*pb.Trainer, error) {
	client := MongoClient

	newTrainers := make([]*models.Trainers, len(request_trainers))
	for i, pbTrainers := range request_trainers {
		newTrainers[i] = MapPbTrainersToModelTrainers(pbTrainers)
	}
	utils.Logger.Info("Trainers to add", zap.Any("newTrainerss", newTrainers))

	var addedTrainers []*pb.Trainer
	for _, trainers := range newTrainers {
		result, err := client.Database("data").Collection("trainers").InsertOne(ctx, trainers)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding value to database")
		}
		objectID, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			trainers.Id = objectID.Hex()
		}
		pbTrainer := MapModelTrainersToPb(trainers)
		addedTrainers = append(addedTrainers, pbTrainer)
	}
	return addedTrainers, nil
}

func GetTrainersFromDb(ctx context.Context, sortOptions primitive.D, filter primitive.M) ([]*pb.Trainer, error) {
	client := MongoClient
	coll := client.Database("data").Collection("trainers")

	var cursor *mongo.Cursor
	var err error
	if len(sortOptions) < 1 {
		cursor, err = coll.Find(ctx, filter)
	} else {
		cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			utils.Logger.Error("Failed to close the cursor", zap.Error(err))
		}
	}()

	trainers, err := decodeEntities(ctx, cursor, func() *pb.Trainer { return &pb.Trainer{} }, newModel)
	if err != nil {
		return nil, err
	}
	return trainers, nil
}

func newModel() *models.Trainers {
	return &models.Trainers{}
}

func UpdateTrainersInDB(ctx context.Context, pbTrainers []*pb.Trainer) ([]*pb.Trainer, error) {
	client := MongoClient

	var updatedTrainers []*pb.Trainer
	for _, trainer := range pbTrainers {

		modelTrainer := MapPbTrainersToModelTrainers(trainer)

		objId, err := primitive.ObjectIDFromHex(trainer.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid Id")
		}

		modelDoc, err := bson.Marshal(modelTrainer)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal error")
		}

		var updateDoc bson.M
		err = bson.Unmarshal(modelDoc, &updateDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal error")
		}

		delete(updateDoc, "_id")

		_, err = client.Database("data").Collection("trainers").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintln("error updating teacher id:", trainer.Id))
		}

		updatedTrainer := MapModelTrainersToPb(modelTrainer)
		updatedTrainers = append(updatedTrainers, updatedTrainer)

	}
	return updatedTrainers, nil
}
