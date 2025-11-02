package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	pb "grpcserver/proto/generated_files"
	"reflect"
)

func decodeEntities[T any, M any](ctx context.Context, cursor *mongo.Cursor, newEntity func() *T, newModel func() *M) ([]*T, error) {
	var entities []*T
	for cursor.Next(ctx) {
		model := newModel()
		err := cursor.Decode(&model)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal Error")
		}
		entity := newEntity()
		modelVal := reflect.ValueOf(model).Elem()
		pbVal := reflect.ValueOf(entity).Elem()
		for i := 0; i < modelVal.NumField(); i++ {
			modelField := modelVal.Field(i)
			modelFieldName := modelVal.Type().Field(i).Name

			pbField := pbVal.FieldByName(modelFieldName)
			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}
		entities = append(entities, entity)
	}
	err := cursor.Err()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal error")
	}
	return entities, nil
}

func MapModelToPb[P any, M any](model M, newPb func() *P) *P {
	pbStruct := newPb()
	modelVal := reflect.Indirect(reflect.ValueOf(model))
	pbVal := reflect.Indirect(reflect.ValueOf(pbStruct))

	for i := 0; i < modelVal.NumField(); i++ {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)
		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbStruct
}

func MapPbToModel[P any, M any](pbStruct P, newModel func() *M) *M {
	modelStruct := newModel()
	pbVal := reflect.Indirect(reflect.ValueOf(pbStruct))
	modelVal := reflect.Indirect(reflect.ValueOf(modelStruct))

	for i := 0; i < pbVal.NumField(); i++ {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return modelStruct
}

func MapModelTrainersToPb(trainers *models.Trainers) *pb.Trainer {
	return MapModelToPb(trainers, func() *pb.Trainer { return &pb.Trainer{} })
}

func MapPbTrainersToModelTrainers(pbTrainer *pb.Trainer) *models.Trainers {
	return MapPbToModel(pbTrainer, func() *models.Trainers { return &models.Trainers{} })
}
