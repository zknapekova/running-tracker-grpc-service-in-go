package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	pb "grpcserver/proto/generated_files"
	"reflect"
)

func AddTrainersToDB(ctx context.Context, request_trainers []*pb.Trainer) ([]*pb.Trainer, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	newTrainers := make([]*models.Trainers, len(request_trainers))
	for i, pbTrainers := range request_trainers {
		newTrainers[i] = MapPbTrainersToModelTrainers(pbTrainers)
	}
	fmt.Println(newTrainers)

	var addedTrainers []*pb.Trainer
	for _, trainers := range newTrainers {
		fmt.Printf("Inserting trainers: %+v\n", trainers)
		result, err := client.Database("main").Collection("trainers").InsertOne(ctx, trainers)
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

func MapModelTrainersToPb(trainers *models.Trainers) *pb.Trainer {
	pbTrainer := &pb.Trainer{}
	modelVal := reflect.Indirect(reflect.ValueOf(trainers))
	pbVal := reflect.Indirect(reflect.ValueOf(pbTrainer))

	for i := 0; i < modelVal.NumField(); i++ {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)
		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbTrainer
}

func MapPbTrainersToModelTrainers(pbTrainer *pb.Trainer) *models.Trainers {
	modelTrainers := models.Trainers{}
	pbVal := reflect.Indirect(reflect.ValueOf(pbTrainer))
	modelVal := reflect.Indirect(reflect.ValueOf(&modelTrainers))

	for i := 0; i < pbVal.NumField(); i++ {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return &modelTrainers
}

func GetTrainersFromDb(ctx context.Context, sortOptions primitive.D, filter primitive.M) ([]*pb.Trainer, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer client.Disconnect(ctx)

	coll := client.Database("main").Collection("trainers")
	var cursor *mongo.Cursor
	if len(sortOptions) < 1 {
		cursor, err = coll.Find(ctx, filter)
	} else {
		cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer cursor.Close(ctx)

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
	defer client.Disconnect(ctx)

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
			return nil, utils.ErrorHandler(err, "internal error")
		}

		delete(updateDoc, "_id")

		_, err = client.Database("main").Collection("trainers").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintln("error updating teacher id:", trainer.Id))
		}

		updatedTrainer := MapModelTrainersToPb(modelTrainer)
		updatedTrainers = append(updatedTrainers, updatedTrainer)

	}
	return updatedTrainers, nil
}
