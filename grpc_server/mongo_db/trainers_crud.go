package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	pb "grpcserver/proto/generated_files"
)

func AddTrainersToDB(ctx context.Context, request_trainers []*pb.Trainer) ([]*pb.Trainer, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer DisconnectMongoClient(client, ctx)

	newTrainers := make([]*models.Trainers, len(request_trainers))
	for i, pbTrainers := range request_trainers {
		newTrainers[i] = MapPbTrainersToModelTrainers(pbTrainers)
	}
	utils.InfoLogger.Println(newTrainers)

	var addedTrainers []*pb.Trainer
	for _, trainers := range newTrainers {
		fmt.Printf("Inserting trainers: %+v\n", trainers)
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
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer DisconnectMongoClient(client, ctx)

	coll := client.Database("data").Collection("trainers")
	var cursor *mongo.Cursor
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
			utils.ErrorLogger.Printf("Failed to close the cursor: %v", err)
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
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer DisconnectMongoClient(client, ctx)

	var updatedTrainers []*pb.Trainer
	for _, trainer := range pbTrainers {
		if trainer.Id == "" {
			return nil, utils.ErrorHandler(errors.New("id cannot be blank"), "id cannot be blank")
		}

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
			return nil, utils.ErrorHandler(err, "internal error")
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
