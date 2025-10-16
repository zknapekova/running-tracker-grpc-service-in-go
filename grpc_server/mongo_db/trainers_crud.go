package mongodb

import (
		"context"
		"fmt"
		 pb "grpcserver/proto/generated_files"
		"grpcserver/internals/models"
		"grpcserver/internals/utils"
		"go.mongodb.org/mongo-driver/bson/primitive"
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
		newTrainers[i] = mapPbTrainersToModelTrainers(pbTrainers)
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
		pbTrainer := mapModelTrainersToPb(trainers)
		addedTrainers = append(addedTrainers, pbTrainer)
	}
	return addedTrainers, nil
}


func mapModelTrainersToPb(trainers *models.Trainers) *pb.Trainer {
	pbTrainer := &pb.Trainer{}
	modelVal := reflect.Indirect(reflect.ValueOf(trainers))
	pbVal := reflect.Indirect(reflect.ValueOf(pbTrainer))

	for i := 0; i < modelVal.NumField(); i++ {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)
		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid () && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbTrainer
}

func mapPbTrainersToModelTrainers(pbTrainer *pb.Trainer) *models.Trainers {
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