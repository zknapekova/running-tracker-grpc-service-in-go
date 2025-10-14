package handlers

import (
		"context"
		"fmt"
		mongodb "grpcserver/mongo_db"
		 pb "grpcserver/proto/generated_files"
		"grpcserver/internals/models"
		"grpcserver/internals/utils"
		"go.mongodb.org/mongo-driver/bson/primitive"
		"reflect"
)

func (s *Server) AddTrainers(ctx context.Context, req *pb.AddTrainersRequest) (*pb.Response, error) {
	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	newTrainers := make([]*models.Trainers, len(req.GetTrainers()))
	for i, pbTrainers := range req.GetTrainers() {
		modelTrainers := models.Trainers{Brand: "Adidas"}
		pbVal := reflect.Indirect(reflect.ValueOf(pbTrainers))
		modelVal := reflect.Indirect(reflect.ValueOf(&modelTrainers))

		for i := 0; i < pbVal.NumField(); i++ {
			pbField := pbVal.Field(i)
			fieldName := pbVal.Type().Field(i).Name

			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(pbField)
			}
		}
		newTrainers[i] = &modelTrainers
	}
	fmt.Println(newTrainers)

	//var addedTrainers []*pb.Trainer
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
		/*pbTrainers := &pb.Trainer{}
		modelVal := reflect.Indirect(reflect.ValueOf(&modelTrainers))
		pbVal := reflect.Indirect(reflect.ValueOf(pbTrainers))

		for i := 0; i < modelVal.NumField(); i++ {
			modelField := modelVal.Field(i)
			modelFieldType := modelVal.Type().Field(i)
			pbField := pbVal.FieldByName(modelFieldType.Name)
			if pbField.IsValid () && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}
		addedTrainers = append(addedTrainers, pbTrainers)*/
	}
	return &pb.Response{
		Message: "Trainers were added to database",
		Code:    0,
	}, nil
}